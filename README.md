# meliChallenge
Meli Challenge

## instalacion
 - Instalar Docker https://www.docker.com/get-started/
 - Instalar Postgresql https://www.postgresql.org/download/
 - Renombrar el archivo .env.example a .env
 - Crear una base de datos nueva en postgres
 - Configurar los valores de la base de datos en el archivo .env (variables que inician en DB_)
 - Configurar el puerto por el que se escucharán las peticiones http en el archivo .env (PORT)
 - Correr el comando make start

## Descripcion general del funcionamiento
Al iniciarse la app se crea el primer climate_record_job, este se encarga de crear los climate_record iniciales (10 años)
También se encarga de guardar contar cuantos climas hubo de cada uno y cual fue el mayor perimetro para los climas de lluvia para guardar esta información en el archivo results que se guarda en el root del proyecto

lluvia:1249 maxima:288
sequia:61
optimo:122

Se parte de un plano cartesiano en el que el sol está siempre hubicado en el origen para facilitar los calculos

Para calcular si hubo un día de sequía se calculan las coordenadas de los planetas según los días que han pasado, se traza una ecuación de la recta entre el planeta mas cercano y el mas lejano al sol, y se calcula si el planeta del medio y el sol está na menos de 0.01km (alignmentThreshold error permitido) de esta linea. 

Para calcular si hubo un día de condiciones optimas se utiliza el metodo descrito anteriormente pero sin tomar en cuenta el sol

Para calcular si hubo un clima con lluvia se calcula con el algoritmo descrito en este video https://www.youtube.com/watch?v=WaYS1gEXEFE&t=542s&ab_channel=huse360

## Descripcion modelos

 - climate_records
   la columna day representa el número de días pasados desde las condiciones iniciales de los planetas (descritos en las variables de entorno)
   la columna climate representa el clima calculado para este día
 - climate_record_jobs
   La aplicación crea los climate_records con un job que se ejecuta por primera vez al ser iniciada que se vuelve a ejecutar cada día a las 00:00 
   En esta tabla se guarda información de los días calculados en cada ejecución
   first_day representa el day del primer climate_record calculado en la ejecucion de este job
   last_day representa el day del ultimo climate_record calculado en la ejecucion de este job
   la columna created_at es la fecha en la que se ejecutó el job, esta se utiliza para calcular cuantos días han pasado desde la ultima ejecucion
   del job para con esto y la columna last_day calcular cuantos días mas se deben calcular.

## Descripcion paquetes

Los paquetes están desacoplados entre ellos y se hacen funcionar mediante inyección de dependencias

 - infrastructure contienela implementacion de la conexión a la base de datos utilizando gorm
 - repository se encarga de persistir los datos utilizando una instancia de la conexión
 - service se encarga de agrupar la logica de cada caso de uso, estos utilizan la lógica de domain para hacer calculos y los paquetes repository para persistir la informacion
 - handler se encarga de guardar los metodos utilizados en las peticiones http, este llama al servicio correspondiente
 - domain
   contiene metodos, estructuras y constantes que pueden ser utilizadas en el resto de paquetes, lógica de dominio y errores pre definidos.
   cuenta con dos subpaquetes meteorology y geometry
 - - geometry
     paquete encargado de recopilar estructuras y metodos necesarios con la geometría del problema (coordenadas, calculo de perimetros, distancia a la recta, si un punto está dentro de un triangulo, etc). este es utilizado en sub paquete meteorology de domain
 - - meteorology
     paquete encargado de recopilar estructuras de planetas y metodos nesarios para calcular el clima, proveer los planetas, calcuar el perimetro entre planetas, etc, está basado en el subpaquete geometry para todos estos calculos
 - app se encarga de instanciar cada paquete, inyectar y agrupar dependencias, definir las rutas, iniciar los jobs, iniciar el servidor, etc

