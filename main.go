package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yama-koo/ent-example/ent"
	"github.com/yama-koo/ent-example/ent/car"
	"github.com/yama-koo/ent-example/ent/group"
	"github.com/yama-koo/ent-example/ent/user"
)

func main() {
	// client, err := ent.Open("mysql", "root:password@tcp(127.0.0.1:3306)/ent?parseTime=true&loc=Local", ent.Debug())
	client, err := ent.Open("mysql", "root:password@tcp(127.0.0.1:3306)/ent?parseTime=true&loc=Local")
	if err != nil {
		log.Fatalf("failed to connect database %+v\n", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalln("failed to creating schema")
	}
	// if err := client.Schema.WriteTo(ctx, os.Stdout); err != nil {
	// 	log.Fatalf("failed printing schema changes: %v", err)
	// }

	// u, err := CreateUser(ctx, client)
	// err = CreateGraph(ctx, client)
	// u, err := QueryUser(ctx, client)
	// _ = QueryCars(ctx, u)
	// _ = QueryCarUsers(ctx, u)
	// _ = QueryGithub(ctx, client)
	// _ = QueryArielCars(ctx, client)
	_ = QueryGroupWithUsers(ctx, client)
	// u, err := CreateCars(ctx, client)
	// if err != nil {
	// 	return
	// }
	// fmt.Printf("%+v\n", u)
}

// func CreateUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
// 	u, err := client.User.Create().SetAge(20).SetName("ディオ").Save(ctx)
// 	if err != nil {
// 		log.Println(err)
// 		return nil, err
// 	}
// 	return u, nil
// }

// func CreateGraph(ctx context.Context, client *ent.Client) error {
// 	// First, create the users.
// 	a8m, err := client.User.
// 		Create().
// 		SetAge(30).
// 		SetName("Ariel").
// 		Save(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	neta, err := client.User.
// 		Create().
// 		SetAge(28).
// 		SetName("Neta").
// 		Save(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	// Then, create the cars, and attach them to the users in the creation.
// 	_, err = client.Car.
// 		Create().
// 		SetModel("Tesla").
// 		SetRegisteredAt(time.Now()). // ignore the time in the graph.
// 		SetOwner(a8m).               // attach this graph to Ariel.
// 		Save(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = client.Car.
// 		Create().
// 		SetModel("Mazda").
// 		SetRegisteredAt(time.Now()). // ignore the time in the graph.
// 		SetOwner(a8m).               // attach this graph to Ariel.
// 		Save(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = client.Car.
// 		Create().
// 		SetModel("Ford").
// 		SetRegisteredAt(time.Now()). // ignore the time in the graph.
// 		SetOwner(neta).              // attach this graph to Neta.
// 		Save(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	// Create the groups, and add their users in the creation.
// 	_, err = client.Group.
// 		Create().
// 		SetName("GitLab").
// 		AddUsers(neta, a8m).
// 		Save(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = client.Group.
// 		Create().
// 		SetName("GitHub").
// 		AddUsers(a8m).
// 		Save(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	log.Println("The graph was created successfully")
// 	return nil
// }

func QueryUser(ctx context.Context, client *ent.Client) (*ent.User, error) {
	u, err := client.User.Query().Where(user.NameEQ("a8m")).Only(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return u, nil
}

func QueryCars(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %v", err)
	}
	log.Println("returned cars:", cars)

	ford, err := a8m.QueryCars().Where(car.ModelEQ("ford")).Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %v", err)
	}
	log.Println(ford)
	return nil
}

func QueryCarUsers(ctx context.Context, a8m *ent.User) error {
	cars, err := a8m.QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed querying user cars: %v", err)
	}

	for _, ca := range cars {
		owner, err := ca.QueryOwner().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying car %q owner: %v", ca.Model, err)
		}
		log.Printf("car %q owner: %q\n", ca.Model, owner.Name)
	}
	return nil
}

func QueryGithub(ctx context.Context, client *ent.Client) error {
	cars, err := client.Group.Query().Where(group.NameEQ("GitHub")).QueryUsers().QueryCars().All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting cars: %v", err)
	}
	log.Println("cars returned:", cars)
	return nil
}

func QueryArielCars(ctx context.Context, client *ent.Client) error {
	// Get "Ariel" from previous steps.
	a8m := client.User.
		Query().
		Where(
			user.HasCars(),
			user.Name("Ariel"),
		).
		OnlyX(ctx)
	cars, err := a8m. // Get the groups, that a8m is connected to:
				QueryGroups(). // (Group(Name=GitHub), Group(Name=GitLab),)
				QueryUsers().  // (User(Name=Ariel, Age=30), User(Name=Neta, Age=28),)
				QueryCars().   //
				Where(car.Not(car.ModelEQ("Mazda"))).
				All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting cars: %v", err)
	}
	log.Println("cars returned:", cars)
	return nil
}

func QueryGroupWithUsers(ctx context.Context, client *ent.Client) error {
	groups, err := client.Group.
		Query().
		Where(group.HasUsers()).
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed getting groups: %v", err)
	}
	log.Println("groups returned:", groups)
	// Output: (Group(Name=GitHub), Group(Name=GitLab),)
	return nil
}

// func CreateCars(ctx context.Context, client *ent.Client) (*ent.User, error) {
// 	tesla, err := client.Car.Create().SetModel("tesla").SetRegisteredAt(time.Now()).Save(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed creating car: %v", err)
// 	}

// 	ford, err := client.Car.Create().SetModel("ford").SetRegisteredAt(time.Now()).Save(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed creating car: %v", err)
// 	}
// 	log.Println("car was created: ", ford)

// 	a8m, err := client.User.Create().SetAge(30).SetName("a8m").AddCars(tesla, ford).Save(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed creating user: %v", err)
// 	}

// 	return a8m, nil
// }
