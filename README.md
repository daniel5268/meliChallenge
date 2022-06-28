# meliChallenge
Meli Challenge

## Instalacion
 - Instalar Docker https://www.docker.com/get-started/
 - Instalar Postgresql https://www.postgresql.org/download/
 - Clonar el repositorio
 - Renombrar el archivo .env.example a .env
 - Crear una base de datos nueva en postgres
 - Configurar los valores de la base de datos en el archivo .env (variables que inician en DB_)
   En macOS docker no puede acceder a hosts externos por lo que la base de datos debe estar en el contenedor y la variable de entorno DB_HOST debe tener el valor host.docker.internal (https://docs.docker.com/desktop/mac/networking/)
 - Configurar el puerto por el que se escucharán las peticiones http en el archivo .env (PORT)
 - Correr el comando `make start`
 - Si se ejecuta el contenedor docker desde macOS se recomienda utilizar el comando anterior sólo para correr las migraciones, por ende finalizar el proceso (control + c), cambiar el valor de DB_HOST de host.docker.internal al nuevo (posible mente localhost) y correr el comando `go run ./src/main.go`

## Descripcion general del funcionamiento
Al iniciarse la app se crea el primer climate_record_job, este se encarga de crear los climate_record iniciales (10 años).
Hay un cron job que se ejecuta todas las noches a las 00:00 que ejecuta otro job para crear más climate_records, esto para asegurar que siempre esten disponibles 10 años a partir de la fecha en la que se consulte, aún así se debe tener en cuenta que los días no se reinician no se borran los dias pasados.

Se parte de un plano cartesiano en el que el sol está siempre hubicado en el origen para facilitar los calculos

Para calcular si hubo un día de sequía se calculan las coordenadas de los planetas según los días que han pasado, se traza una ecuación de la recta entre el planeta mas cercano y el mas lejano al sol, y se calcula si el planeta del medio y el sol está na menos de 0.01km (alignmentThreshold error permitido) de esta linea. 

Para calcular si hubo un día de condiciones optimas se utiliza el metodo descrito anteriormente pero sin tomar en cuenta el sol

Para calcular si hubo un clima con lluvia se calcula con el algoritmo descrito en este video https://www.youtube.com/watch?v=WaYS1gEXEFE&t=542s&ab_channel=huse360

## Descripcion modelos

 - climate_records
   la columna day representa el número de días pasados desde las condiciones iniciales de los planetas (descritos en las variables de entorno)
   la columna climate representa el clima calculado para este día
   la columna perimeter solo se calcula si el clima es de lluvia y guarda cual es el perimetro entre los planetas para calcular el dia con mayor intensidad de lluvia
 - climate_record_jobs
   La aplicación crea los climate_records con un job que se ejecuta por primera vez al ser iniciada, se vuelve a ejecutar cada día a las 00:00, este job agrega un clima por cada que ha pasado desde la ultima ejecución, si se ejecuta efectivamente todos días cada día agregará sólo un climate_record
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

## Endpoints

Para obtener el clima de un día específico
curl --location --request GET 'https://melimeteorology.herokuapp.com/api/challenge/clima?dia=566'

Para obtener cuantos días hubo de cada clima y cuál fue el día de lluvia maxima
curl --location --request GET 'https://melimeteorology.herokuapp.com/api/challenge/clima/resumen?primer_dia=0&ultimo_dia=3650'

## Errores
El servidor cuenta con un middleware de errores para formatear a una respuesta amigable al usuario los posibles errores
