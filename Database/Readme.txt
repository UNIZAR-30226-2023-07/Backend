No instaleis docker en windows -> crear una maquina virtual con docker 
(en caso de que la maquina virtual os vaya muy lenta (mi caso) conectar con 
una red anfitrión la maquina virtual y así al ejecutar docker podeis ejecutar 
pgAdmin en un navegador que no sea el de la maquina virtual)

Para crear y ejecutar la base de datos -> docker-compose up

Para entrar en pgAdmin4 entrar escribir en un navegador localhost
El usuario de pgAdmin4 es frances@allen.es y contraseña 1234
Una vez dentro habra que añadir el servidor, para ello click derecho en server y pulsar en register-> server
Escribir en name Pro_Soft y en address postgres y en contraseña 1234 y save y aparecerá la base de datos

Una vez hecho esto se puede iniciar desde la aplicación de docker

Ejecutar el main para ejecutar las pruebas


