package main

import (
	"log"
	"net/http"
	"github.com/gofiber/fiber/v2"
	// "github.com/joho/godotenv"
	"os"
	"github.com/RionDsilvaCS/kkhomeneeds/models"
	"github.com/RionDsilvaCS/kkhomeneeds/storage"
	"database/sql"

	jwtware "github.com/gofiber/contrib/jwt"
    "github.com/golang-jwt/jwt/v5"
	"time"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
} 

type Repository struct {
	DB *sql.DB
}

func(r *Repository) CreateProduct(context *fiber.Ctx) error{
	product := models.Products{} 

	err := context.BodyParser(&product)

	if err != nil {
		return err
	}

	return nil

}

func(r *Repository) GetProducts(context *fiber.Ctx) error{
	productModels := []models.Products{} 

	rows, err := r.DB.Query("select * from products")
	if err != nil{
		log.Fatal(err)
	}
	for rows.Next(){
		var t models.Products
		err := rows.Scan(&t.ID, &t.Img_1, &t.Img_2, &t.Title, &t.Description, &t.Status, &t.MRP_price, &t.Discount_price)
		if err != nil{
			log.Fatal(err)
		}
		productModels = append(productModels, t)
	}


	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "products fetched successfully",
		"data": productModels,
	})

	return nil
}

func(r *Repository) UserLogin(context *fiber.Ctx) error{
	credentials := Credential{}
    
	if err := context.BodyParser(&credentials); err != nil {
		return err
	}

    if credentials.Username != "john" || credentials.Password != "doe" {
        return context.SendStatus(fiber.StatusUnauthorized)
    }

    // Create the Claims
    claims := jwt.MapClaims{
        "name":  "a",
        "admin": false,
        "exp":   time.Now().Add(time.Hour * 72).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    t, err := token.SignedString([]byte("secret"))
    if err != nil {
        return context.SendStatus(fiber.StatusInternalServerError)
    }

    return context.JSON(fiber.Map{"token": t})
}

func(r *Repository) Restricted(context *fiber.Ctx) error{
	user := context.Locals("user").(*jwt.Token)
    claims := user.Claims.(jwt.MapClaims)
    name := claims["name"].(string)
    return context.SendString("Welcome " + name)
}

func(r *Repository) SetUpRoutes(app *fiber.App){

	api_product := app.Group("/product")

	api_product.Get("/products", r.GetProducts)
	// api.Get("/get_products/:id", r.GetProductByID)
	api_product.Post("/create_products", r.CreateProduct)
	// api.Delete("/delete_product", r.DeleteProduct)


	api_user := app.Group("/user")

	api_user.Post("/login", r.UserLogin)
	api_user.Use(jwtware.New(jwtware.Config{
        SigningKey: jwtware.SigningKey{Key: []byte("secret")},
    }))
	api_user.Get("/afterlogin", r.Restricted)
}



func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	config := &storage.Config{
		Host: os.Getenv("DB_HOST"),
		Password: os.Getenv("DB_PASS"),
		User: os.Getenv("DB_USER"),
		DBName: os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not connect")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetUpRoutes(app)

	// app.Listen(":3000")

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}
	
	log.Fatal(app.Listen("0.0.0.0:" + port))

}