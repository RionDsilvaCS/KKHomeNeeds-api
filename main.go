package main

import (
	"log"
	"net/http"
	"github.com/gofiber/fiber/v2"
	"os"
	"github.com/RionDsilvaCS/kkhomeneeds/models"
	"github.com/RionDsilvaCS/kkhomeneeds/storage"
	"database/sql"

	jwtware "github.com/gofiber/contrib/jwt"
    "github.com/golang-jwt/jwt/v5"
	"time"

	"github.com/joho/godotenv"
	"fmt"
)

type Credential struct {
	Email string `json:"email"`
	Password string `json:"password"`
} 

type User struct {
	UserID 	int 			`json:"userid"`
	Username string			`json:"username"`
	Email string			`json:"email"`
	Password string			`json:"password"`
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

func(r *Repository) UserSignin(context *fiber.Ctx) error{
	userDetails := User{}

	if err := context.BodyParser(&userDetails); err != nil {
		return err
	}
	query := fmt.Sprintf("INSERT INTO USERS (ID, USERNAME, EMAIL, PASSCODE) SELECT %d, '%s', '%s', '%s' WHERE NOT EXISTS (SELECT 1 FROM USERS WHERE EMAIL = '%s');", userDetails.UserID, userDetails.Username, userDetails.Email, userDetails.Password, userDetails.Email)
	var output string
	err := r.DB.QueryRow(query).Scan(&output)
	if err != nil {
		return context.JSON(fiber.Map{"message": "Account already exits", "output": err})
	}

	return context.JSON(fiber.Map{"message": "New Account created", "output": output})
}


func(r *Repository) UserLogin(context *fiber.Ctx) error{
	credentials := Credential{}
    
	if err := context.BodyParser(&credentials); err != nil {
		return err
	}
	user := User{}
	query := fmt.Sprintf("SELECT USERNAME FROM USERS WHERE EMAIL = '%s' AND PASSCODE = '%s';", credentials.Email, credentials.Password)
	err := r.DB.QueryRow(query).Scan(&user.Username)
	if err != nil {
		return context.SendStatus(fiber.StatusUnauthorized)
	}

    claims := jwt.MapClaims{
        "name":  user.Username,
        "admin": false,
        "exp":   time.Now().Add(time.Hour * 72).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    t, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET"))) 
    if err != nil {
        return context.SendStatus(fiber.StatusInternalServerError)
    }

    return context.JSON(fiber.Map{"token": t, "username": user.Username})
}

func(r *Repository) Restricted(context *fiber.Ctx) error{
	user := context.Locals("user").(*jwt.Token)
    claims := user.Claims.(jwt.MapClaims)

    name := claims["name"].(string)
    return context.SendString("Welcome " + name)
}




func(r *Repository) SetUpRoutes(app *fiber.App){

	api_product := app.Group("/product")

	api_product.Use(jwtware.New(jwtware.Config{
        SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("TOKEN_SECRET"))},
    }))
	api_product.Get("/products", r.GetProducts)
	// api.Get("/get_products/:id", r.GetProductByID)
	api_product.Post("/create_products", r.CreateProduct)
	// api.Delete("/delete_product", r.DeleteProduct)


	api_user := app.Group("/user")

	api_user.Post("/signin", r.UserSignin)
	api_user.Post("/login", r.UserLogin)

	api_user.Use(jwtware.New(jwtware.Config{
        SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("TOKEN_SECRET"))},
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