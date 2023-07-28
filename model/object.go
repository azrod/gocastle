// models/objects.go

package model

import "fmt"

// ObjectStat represents a specific stat of an object.
type ObjectStat struct {
	Name     string // Stat name, e.g., "Strength", "Health", etc.
	Modifier int    // Modifier value for the stat.
}

// Object represents an object with its properties.
type Object struct {
	Name     string       // Object name.
	Category string       // Object category.
	Weight   int          // Object weight in grams
	Equipped bool         // Is object equipped
	Stats    []ObjectStat // Object stats (e.g., strength, health, etc.).
}

// Category represents a common object category.
type Category struct {
	Name        string // Category name.
	Description string // Description of the category.
}

var (
	CategoryList []Category

	HuntingKnife = Object{
		Name:     "Hunting knife",
		Category: "Weapon",
		Weight:   200,
		Stats: []ObjectStat{
			{
				Name:     "physicalDamage",
				Modifier: 2,
			},
		},
	}

	BareHands = Object{
		Name:     "Bare Hands",
		Category: "Weapon",
		Weight:   0,
	}
)

// GenerateCategories creates all the categories based on the provided names and descriptions.
func GenerateCategories(names, descriptions []string) {
	if len(names) != len(descriptions) {
		panic("number of names and descriptions should be the same")
	}

	CategoryList = make([]Category, len(names))
	for i := range names {
		CategoryList[i] = Category{
			Name:        names[i],
			Description: descriptions[i],
		}
	}
}

// InitializeCategories initializes the categories with their names and descriptions.
func InitializeCategories() {
	GenerateCategories(
		[]string{
			"Weapon",
			"Body Armor",
			"Head Gear",
			"Belt Item",
			"Neckwear",
			"Overgarment",
			"Shield",
			"Gauntlets",
			"Boots",
			"Purse",
			"Belt",
			"Left Ring",
			"Right Ring",
			"Bracers",
		},
		[]string{
			"Weapons used for combat.",
			"Gear worn to the chest.",
			"Head gear (can be hats, helmets,...).",
			"Consumables that are easily accessible in combat.",
			"Items worn around the neck.",
			"Outer garments like cloaks or capes.",
			"Shields used for defense.",
			"Protective gloves for hands.",
			"Footwear.",
			"A container for money.",
			"A belt worn around the waist.",
			"A ring worn on the left hand.",
			"A ring worn on the right hand.",
			"Arm protectors.",
		},
	)
}

// CategoryExists checks if the given category exists in the CommonCategories slice.
func CategoryExists(categoryName string) bool {
	for _, cat := range CategoryList {
		if cat.Name == categoryName {
			return true
		}
	}
	return false
}

// CreateObject creates a copy of the given object and returns it.
// It also validates the category before creating the object.
func CreateObject(obj Object) (Object, error) {
	// Validate the category.
	if !CategoryExists(obj.Category) {
		return Object{}, fmt.Errorf("category '%s' does not exist", obj.Category)
	}

	// Create a new Object with the same properties.
	newObject := Object{
		Name:     obj.Name,
		Category: obj.Category,
		Weight:   obj.Weight,
	}

	// Copy the ObjectStat slice.
	for _, stat := range obj.Stats {
		newObject.Stats = append(newObject.Stats, ObjectStat{
			Name:     stat.Name,
			Modifier: stat.Modifier,
		})
	}

	return newObject, nil
}
