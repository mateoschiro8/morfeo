#pagebreak()

= Introducción teórica
El objetivo de este trabajo es la creación de un sistema funcional para el armado y alerta de *honeytokens* (señuelo digital que funciona como sistema de alarma; al accederlo, se delata la presencia de una actividad no autorizada en el sistema y una posible brecha de seguridad). 

Existen muchos tipos de *honeytokens*, pueden ser archivos PDF, Word/Excel, binarios, Direcciones de correo, Cuentas de usuario falsas entre otras. En este trabajo, investigamos 
los mecanismos _Call Home_, inspirados en herramientas como _CanaryTokens_, para desarrollar *honeytokens* en los siguientes formatos: QR, Binarios, PDF, IMG y CSS.

- QR: se encapsula una URL configurada para actuar como nuestro canary token. Al ser escaneado por un atacante, se genera una solicitud hacia el servidor que controla el investigador. Al recibir esta solicitud en un endpoint, se constituye el evento de callback y la confirmación inmediata de la intrusión o el acceso no autorizado. El servidor luego activa la alerta de seguridad sin necesidad de interactuar o exponer datos sensibles y se redirige al usuario a otra URL con un codigo 302.

- Binario: como el caso anterior, se hace un callback a nuestro endpoint controlado pero se activa de otra manera. En este *honeytoken* hacemos un wrapper de nuestro codigo de alerta sobre un programa legitimo. Cuando se ejecuta el binario original se activa nuestra alerta que llega al servidor mediante el pedido HTTP mientras que el programa original corre sin problema.

- PDF: en este tipo de *honeytoken* la activacion de la alarma se efectua cuando se abre el documento en cuestion. Una de las particularidades es que hay variantes respecto a como puede hacerse esto, ya sea con codigo JavaScript embebebido que realiza el pedido HTTP o con la busqueda a un recurso externo (como una referencia a una imagen que apunta a una URL alojada nuestro servidor). El pedido al servidor se hace con el lector, que ejecuta el script de inicio, y el codigo JS incrustado realiza la petición HTTP o este interpreta la referencia externa, intentando descargar la imagen para mostrar el PDF y termina enviando la solicitud al servidor.

- CSS: 