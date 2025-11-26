#v(2cm)

#align(center, [
    #text(size: 32pt, weight: "bold")[TP1: TokenSnare]
    
    #v(0.8cm)

    #text(size: 22pt, weight: "semibold")[Seguridad de la información]

    #text(size: 18pt, weight: "semibold")[2do cuatrimestre 2025]
    
    #v(1cm)

    #image("img/honeypot.jpg", width: 40%)

    #v(1cm)

    #text(size: 18pt, weight: "semibold")[Ciberseguros]

    #set table(
        fill: (_, y) => if y == 0 { rgb("EAF2F5") },
    )

    #table(
        columns: 3,
        table.header[*Nombre*][*LU*][*Mail*],
        [Juan Begalli],[139/22 ],[juanbegalli\@gmail.com],
        [Francisco Cueto],[223/22],[francue3\@gmail.com],
        [Santiago Rivas],[415/22],[santiagorivas0203\@gmail.com],
        [Mateo Schiro],[657/22],[mateo.schiro8\@gmail.com],
    )   
])