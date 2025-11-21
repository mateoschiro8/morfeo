package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"seehuhn.de/go/pdf"
)

var (
	inputPath  string
	outputPath string
	imageURL   string = "http://localhost:8000/track"
	uri        string
)

var pdfCmd = &cobra.Command{
	Use:   "pdf",
	Short: "Genera el honeytoken de pdf",
	Run:   run,
}

func init() {
	pdfCmd.Flags().StringVar(&inputPath, "in", "", "Ruta al archivo de entrada")
	pdfCmd.Flags().StringVar(&outputPath, "out", "", "Ruta al archivo de salida")
	pdfCmd.MarkFlagRequired("in")
	pdfCmd.MarkFlagRequired("out")
	rootCmd.AddCommand(pdfCmd)
}

func run(cmd *cobra.Command, args []string) {

	// ---------------------------
	// 1) Reader del PDF original
	// ---------------------------
	inFile, err := os.Open(inputPath)
	if err != nil {
		panic(fmt.Errorf("error abriendo input: %w", err))
	}
	defer inFile.Close()

	r, err := pdf.NewReader(inFile, nil)
	if err != nil {
		panic(fmt.Errorf("error creando reader: %w", err))
	}

	// ---------------------------
	// 2) Writer del PDF nuevo
	// ---------------------------
	outFile, err := os.Create(outputPath)
	if err != nil {
		panic(fmt.Errorf("error creando output: %w", err))
	}
	defer outFile.Close()

	w, err := pdf.NewWriter(outFile, r.GetMeta().Version, nil)
	if err != nil {
		panic(fmt.Errorf("error creando writer: %w", err))
	}

	// ----------------------------------------------------
	// 3) Obtener lista de páginas recorriendo el árbol PDF
	//    (inline, sin funciones externas)
	// ----------------------------------------------------
	var collectPages func(ref pdf.Reference) ([]pdf.Reference, error)
	collectPages = func(ref pdf.Reference) ([]pdf.Reference, error) {
		obj, err := r.Get(ref, true)
		if err != nil {
			return nil, err
		}
		dict, ok := obj.(pdf.Dict)
		if !ok {
			return nil, fmt.Errorf("nodo Pages/Page inválido")
		}

		kidsObj, err := pdf.Resolve(r, dict["Kids"])
		if err != nil {
			return nil, fmt.Errorf("error resolviendo Kids: %w", err)
		}

		kids, ok := kidsObj.(pdf.Array)
		if !ok {
			return nil, fmt.Errorf("Kids no es array")
		}

		var out []pdf.Reference

		for _, k := range kids {
			kidRef, ok := k.(pdf.Reference)
			if !ok {
				return nil, fmt.Errorf("kid no es referencia")
			}

			kidObj, err := r.Get(kidRef, true)
			if err != nil {
				return nil, err
			}

			kidDict, ok := kidObj.(pdf.Dict)
			if !ok {
				return nil, fmt.Errorf("kid no es dict")
			}

			t, _ := kidDict["Type"].(pdf.Name)

			switch t {
			case "Page":
				out = append(out, kidRef)

			case "Pages":
				nested, err := collectPages(kidRef)
				if err != nil {
					return nil, err
				}
				out = append(out, nested...)

			default:
				return nil, fmt.Errorf("Kid con Type inesperado: %v", t)
			}
		}

		return out, nil
	}

	pages, err := collectPages(r.GetMeta().Catalog.Pages)
	if err != nil {
		panic(fmt.Errorf("error obteniendo páginas: %w", err))
	}

	if len(pages) == 0 {
		panic(fmt.Errorf("el PDF no tiene páginas"))
	}

	// -------------------------------
	// 4) Copiar cada página al writer
	// -------------------------------
	newPages := make([]pdf.Reference, len(pages))

	for i, p := range pages {
		obj, err := r.Get(p, true)
		if err != nil {
			panic(fmt.Errorf("error obteniendo página %d: %w", i, err))
		}

		ref := w.Alloc()
		newPages[i] = ref

		if err := w.Put(ref, obj); err != nil {
			panic(fmt.Errorf("error copiando página %d: %w", i, err))
		}
	}

	// -------------------------------------------------------
	// 5) Modificar la PRIMERA página agregando imagen remota
	// -------------------------------------------------------
	origFirst := pages[0]
	newFirst := newPages[0]

	firstObj, err := r.Get(origFirst, true)
	if err != nil {
		panic(fmt.Errorf("error obteniendo primera página: %w", err))
	}

	pageDict, ok := firstObj.(pdf.Dict)
	if !ok {
		panic(fmt.Errorf("primera página no es un Dict"))
	}

	// --- Resources ---
	resObj := pageDict["Resources"]
	var resDict pdf.Dict

	if resObj == nil {
		resDict = pdf.Dict{}
	} else {
		resolved, _ := pdf.Resolve(r, resObj)
		resDict, _ = resolved.(pdf.Dict)
		if resDict == nil {
			resDict = pdf.Dict{}
		}
	}

	// --- XObject dict ---
	xoObj := resDict["XObject"]
	var xoDict pdf.Dict

	if xoObj != nil {
		resolved, _ := pdf.Resolve(r, xoObj)
		xoDict, _ = resolved.(pdf.Dict)
	}
	if xoDict == nil {
		xoDict = pdf.Dict{}
	}

	// --- Crear XObject "remoto" ---
	imgRef := w.Alloc()
	imgDict := pdf.Dict{
		"Type":             pdf.Name("XObject"),
		"Subtype":          pdf.Name("Image"),
		"Width":            pdf.Integer(1),
		"Height":           pdf.Integer(1),
		"ColorSpace":       pdf.Name("DeviceGray"),
		"BitsPerComponent": pdf.Integer(8),
		"Length":           pdf.Integer(0),
		"F": pdf.Dict{
			"FS": pdf.Name("URL"),
			"F":  pdf.String(imageURL),
		},
	}

	if err := w.Put(imgRef, imgDict); err != nil {
		panic(fmt.Errorf("error escribiendo imagen remota: %w", err))
	}

	xoDict["Img0"] = imgRef
	resDict["XObject"] = xoDict
	pageDict["Resources"] = resDict

	// --- Stream de contenido que dibuja la imagen ---
	contentRef := w.Alloc()
	stream := "q 1 0 0 1 0 0 cm /Img0 Do Q"

	st, err := w.OpenStream(contentRef, pdf.Dict{}, pdf.FilterCompress{})
	if err != nil {
		panic(fmt.Errorf("error creando stream: %w", err))
	}
	_, err = st.Write([]byte(stream))
	if err != nil {
		panic(err)
	}
	if err := st.Close(); err != nil {
		panic(err)
	}

	pageDict["Contents"] = contentRef

	// --- Escribir página modificada ---
	if err := w.Put(newFirst, pageDict); err != nil {
		panic(fmt.Errorf("error escribiendo página modificada: %w", err))
	}

	// -------------------------------
	// 6) Cerrar el writer → FIN PDF
	// -------------------------------
	if err := w.Close(); err != nil {
		panic(fmt.Errorf("error cerrando writer: %w", err))
	}
}

func runWithAnnotations(cmd *cobra.Command, args []string) {

	fmt.Println(os.Getwd())

	in, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		panic(err)
	}

	pdfout, err := os.Open(outputPath)
	if err != nil {
		panic(err)
	}
	defer pdfout.Close()

	r, err := pdf.NewReader(pdfout, nil)
	if err != nil {
		panic(err)
	}

	pagesRef := r.GetMeta().Catalog.Pages
	if pagesRef == 0 {
		panic(err)
	}

	pagesObj, err := r.Get(pagesRef, true)
	if err != nil {
		panic(err)
	}
	pagesDict, ok := pagesObj.(pdf.Dict)
	if !ok {
		panic(err)
	}

	kidsObj, err := pdf.Resolve(r, pagesDict["Kids"])
	if err != nil {
		panic(err)
	}
	kids, ok := kidsObj.(pdf.Array)
	if !ok || len(kids) == 0 {
		panic(err)
	}

	firstPageRef, ok := kids[0].(pdf.Reference)
	if !ok {
		panic(err)
	}

	firstPageObj, err := r.Get(firstPageRef, true)
	if err != nil {
		panic(err)
	}
	pageDict, ok := firstPageObj.(pdf.Dict)
	if !ok {
		panic(err)
	}

	w, err := pdf.NewWriter(out, r.GetMeta().Version, nil)
	if err != nil {
		panic(err)
	}

	annotation := pdf.Dict{
		"Type":    pdf.Name("Annot"),
		"Subtype": pdf.Name("Link"),
		"Rect":    pdf.Array{pdf.Number(0), pdf.Number(0), pdf.Number(0), pdf.Number(0)},
		"Border":  pdf.Array{pdf.Number(0), pdf.Number(0), pdf.Number(0)},
		"A": pdf.Dict{
			"S":   pdf.Name("URI"),
			"URI": pdf.String(uri),
		},
	}

	annotRef := w.Alloc()
	err = w.Put(annotRef, annotation)
	if err != nil {
		panic(err)
	}

	// Obtener o crear el array de anotaciones
	annots := pageDict["Annots"]
	if annots == nil {
		annots = pdf.Array{}
	} else {
		// Resolver si es una referencia
		annots, err = pdf.Resolve(r, annots)
		if err != nil {
			panic(err)
		}
	}

	// Añadir la nueva anotación
	annotsArray, ok := annots.(pdf.Array)
	if !ok {
		annotsArray = pdf.Array{}
	}
	pageDict["Annots"] = append(annotsArray, annotRef)

	// Actualizar la página
	err = w.Put(firstPageRef, pageDict)
	if err != nil {
		panic(err)
	}

	w.Close()

}
