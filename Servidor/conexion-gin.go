package main

import (
	"Handlers"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
	"encoding/json"
)

type Mensaje struct {
    Tipo  string `json:"tipo"`
    Contenido string `json:"contenido"`
}

func main() {
	// Set the router as the default one shipped with Gin
	router := gin.Default()
	m := melody.New()

	router.LoadHTMLFiles("chan.html")
	router.Use(static.Serve("/", static.LocalFile(".", true)))

	// Setup route group for the API
	api := router.Group("/api")
	{
		//Procesa una petición de login
		api.POST("/auth/login", Handlers.PostLogin)

		//Procesa una petición de registro
		api.POST("/auth/register", Handlers.PostRegister)

		//Ejemplo de paso de parametros por url
		api.GET("/prueba/:param", getParam)

		//Ejemplo para devolver más de un dato
		api.GET("/prueba/names", getNames)

		//Ejemplo para devolver structs de datos
		//api.GET("/prueba/users", getUsers)

		//Carga la página del chat/lobby
		api.GET("/channel/:lobby", func(c *gin.Context) {
			c.HTML(http.StatusOK, "chan.html", nil)
		})

		api.GET("/ws/:lobby", func(c *gin.Context) {
			m.HandleRequest(c.Writer, c.Request)
		})

		m.HandleMessage(func(s *melody.Session, msg []byte) {
			var mensaje Mensaje
			err := json.Unmarshal(msg, &mensaje)
			if err != nil {
				fmt.Println("Error al decodificar mensaje:", err)
				return
			} else {
				fmt.Println("Mensaje recibido: ", mensaje.Tipo, mensaje.Contenido)
				// procesar mensaje
				// ...
				// enviar con jsonMsg,err := json.Marshal(mensaje)

			}
			m.BroadcastFilter(msg, func(q *melody.Session) bool { //Envia la información a todos con la misma url
				return q.Request.URL.Path == s.Request.URL.Path
			})
		})

	}

	// Start and run the server
	router.Run(":3001")
}

func getParam(c *gin.Context) {
	param := c.Param("param")
	c.JSON(http.StatusOK, gin.H{
		"name": param,
	})

}

func getNames(c *gin.Context) {
	var param = []string{"Adolfo", "Adrián", "Agustín", "Aitor", "Aitor-tilla"}
	c.JSON(http.StatusOK, gin.H{
		"name": param,
	})

}

/*
func getUsers(c *gin.Context) {
	var users = []Login{{"a@gmail.com", "1234"}, {"b@gmail.com", "5678"}, {"c@gmail.com", "91011"}, {"d@gmail.com", "1234"}}
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
*/
