package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"

	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

var mg MongoInstance

const dbName = "fiber-hrms"

var mongoURI = "mongodb+srv://hrms.rloww82.mongodb.net/?retryWrites=true&w=majority?directConnection=True"

// &w=majority

type User struct {
	ID          string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string  `json:"name"`
	Age         float64 `json:"age"`
	Location    string  `json:"location"`
	Blood_group string  `json:"blood_group"`
	Contact     float64 `json:"contact"`
	Email       string  `json:"email"`
}

func Connect() error {
	credential := options.Credential{
		Username: "Dhruvisha01",
		Password: "Moose_snow@2021",
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI).SetAuth(credential))
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	db := client.Database(dbName)

	if err != nil {
		fmt.Printf("Error in connection!")
		return err
	}

	mg = MongoInstance{
		Client: client,
		Db:     db,
	}
	return nil
}

func main() {
	fmt.Println(mongoURI)
	if err := Connect(); err != nil {
		fmt.Printf("Error part 2")
		log.Fatal(err)
	}
	app := fiber.New()

	app.Get("/users", func(c *fiber.Ctx) error {
		query := bson.D{{}}

		cursor, err := mg.Db.Collection("users").Find(c.Context(), query)

		if err != nil {
			// fmt.Printf("Error part 3")
			return c.Status(500).SendString(err.Error())
		}

		var users []User = make([]User, 0)

		if err := cursor.All(c.Context(), &users); err != nil {
			fmt.Printf("Error 4")
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(users)
	})

	app.Post("/user", func(c *fiber.Ctx) error {
		collection := mg.Db.Collection("users")

		user := new(User)

		if err := c.BodyParser(user); err != nil {
			fmt.Printf("Error 5")
			return c.Status(400).SendString(err.Error())
		}

		user.ID = ""

		insertionResult, err := collection.InsertOne(c.Context(), user)

		if err != nil {
			fmt.Printf("Error g - %v", err)
			return c.Status(500).SendString(err.Error())
		}

		filter := bson.D{{Key: "_id", Value: insertionResult.InsertedID}}
		createdRecord := collection.FindOne(c.Context(), filter)

		createdUser := &User{}

		createdRecord.Decode(createdUser)

		return c.Status(201).JSON(createdUser)

	})

	app.Get("/users/:blood_group/:location", func(c *fiber.Ctx) error {
		blood_group := c.Params("blood_group")
		location := c.Params("location")

		fmt.Println(blood_group)
		fmt.Println(location)

		filter := bson.D{
			{Key: "$and",
				Value: bson.A{
					bson.D{{Key: "location", Value: location}},
					bson.D{
						{Key: "$or",
							Value: bson.A{
								bson.D{{Key: "blood_group", Value: blood_group}},
							},
						},
					},
				}},
		}

		if blood_group == "A+ve" {
			filter = bson.D{
				{Key: "$and",
					Value: bson.A{
						bson.D{{Key: "location", Value: location}},
						bson.D{
							{Key: "$or",
								Value: bson.A{
									bson.D{{Key: "blood_group", Value: blood_group}},
									bson.D{{Key: "blood_group", Value: "O+ve"}},
									bson.D{{Key: "blood_group", Value: "O-ve"}},
									bson.D{{Key: "blood_group", Value: "A-ve"}},
								},
							},
						},
					}},
			}
		} else if blood_group == "AB+ve" {
			filter = bson.D{
				{Key: "$and",
					Value: bson.A{
						bson.D{{Key: "location", Value: location}},
						bson.D{
							{Key: "$or",
								Value: bson.A{
									bson.D{{Key: "blood_group", Value: blood_group}},
									bson.D{{Key: "blood_group", Value: "A+ve"}},
									bson.D{{Key: "blood_group", Value: "B+ve"}},
									bson.D{{Key: "blood_group", Value: "O+ve"}},
									bson.D{{Key: "blood_group", Value: "A-ve"}},
									bson.D{{Key: "blood_group", Value: "B-ve"}},
									bson.D{{Key: "blood_group", Value: "O-ve"}},
									bson.D{{Key: "blood_group", Value: "AB-ve"}},
								},
							},
						},
					}},
			}
		} else if blood_group == "B+ve" {
			filter = bson.D{
				{Key: "$and",
					Value: bson.A{
						bson.D{{Key: "location", Value: location}},
						bson.D{
							{Key: "$or",
								Value: bson.A{
									bson.D{{Key: "blood_group", Value: blood_group}},
									bson.D{{Key: "blood_group", Value: "B-ve"}},
									bson.D{{Key: "blood_group", Value: "O+ve"}},
									bson.D{{Key: "blood_group", Value: "O-ve"}},
								},
							},
						},
					}},
			}
		} else if blood_group == "O+ve" {
			filter = bson.D{
				{Key: "$and",
					Value: bson.A{
						bson.D{{Key: "location", Value: location}},
						bson.D{
							{Key: "$or",
								Value: bson.A{
									bson.D{{Key: "blood_group", Value: blood_group}},
									bson.D{{Key: "blood_group", Value: "O-ve"}},
								},
							},
						},
					}},
			}
		} else if blood_group == "A-ve" {
			filter = bson.D{
				{Key: "$and",
					Value: bson.A{
						bson.D{{Key: "location", Value: location}},
						bson.D{
							{Key: "$or",
								Value: bson.A{
									bson.D{{Key: "blood_group", Value: blood_group}},
									bson.D{{Key: "blood_group", Value: "O-ve"}},
								},
							},
						},
					}},
			}
		} else if blood_group == "AB-ve" {
			filter = bson.D{
				{Key: "$and",
					Value: bson.A{
						bson.D{{Key: "location", Value: location}},
						bson.D{
							{Key: "$or",
								Value: bson.A{
									bson.D{{Key: "blood_group", Value: blood_group}},
									bson.D{{Key: "blood_group", Value: "O-ve"}},
									bson.D{{Key: "blood_group", Value: "A-ve"}},
									bson.D{{Key: "blood_group", Value: "B-ve"}},
								},
							},
						},
					}},
			}
		} else if blood_group == "B-ve" {
			filter = bson.D{
				{Key: "$and",
					Value: bson.A{
						bson.D{{Key: "location", Value: location}},
						bson.D{
							{Key: "$or",
								Value: bson.A{
									bson.D{{Key: "blood_group", Value: blood_group}},
									bson.D{{Key: "blood_group", Value: "O-ve"}},
								},
							},
						},
					}},
			}
		} else if blood_group == "O-ve" {
			filter = bson.D{
				{Key: "$and",
					Value: bson.A{
						bson.D{{Key: "location", Value: location}},
						bson.D{
							{Key: "$or",
								Value: bson.A{
									bson.D{{Key: "blood_group", Value: blood_group}},
								},
							},
						},
					}},
			}
		}

		coll := mg.Db.Collection("users")

		var results []User
		cursor, err := coll.Find(context.TODO(), filter)
		if err != nil {
			if err == mongo.ErrNoDocuments {

				// // This error means your query did not match any documents.
				return c.SendStatus(500)
			}
			panic(err)
		}

		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}

		return c.Status(200).JSON(results)
	})
	log.Fatal(app.Listen(":3000"))
}
