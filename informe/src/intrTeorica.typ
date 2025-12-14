#pagebreak()

= Introducción
El objetivo de este trabajo es la creación de un sistema funcional para el armado y alerta de *honeytokens*. Los honeytokens son artefactos digitales deliberadamente falsos que no tienen uso legítimo, pero están diseñados para parecer reales. Su propósito es detectar accesos no autorizados, exfiltración de datos o movimiento lateral: cualquier interacción con un honeytoken es por definición sospechosa, y genera una señal de alerta. 

Se utilizan tanto en prevención y detección de intrusiones como en análisis forense, ya que permiten identificar vectores de ataque, superficies expuestas y comportamientos del adversario sin interferir con sistemas productivos.  

Existen muchos tipos de honeytokens, pueden ser archivos PDF, Word/Excel, binarios, direcciones de correo, cuentas de usuario falsas, registros DNS, entradas en bases de datos, entre otras. En este trabajo, desarrolamos una herramienta capaz de crear honeytokens en los siguientes formatos: QR, Binarios, PDF, IMG y CSS.

Todos los tokens desarrollados en este trabajo comparten el mismo mecanismo de _Call Home_: un pedido HTTP/HTTPS a un endpoint controlado en nuestro servidor, desde el cual podemos detectar el pedido y realizar la alerta. Esto lo logramos mediante un pedido directo, o el uso de un recurso externo que está alojado en nuestro servidor.

Otros mecanismos de activación pueden ser: DNS callbacks, donde el uso del token provoca una resolución DNS hacia un dominio controlado; SMTP callbacks, típicos en credenciales o direcciones que generan un envío de correo;  integraciones con servicios cloud (por ejemplo, accesos a buckets, funciones o APIs que registran eventos); entre otros.
