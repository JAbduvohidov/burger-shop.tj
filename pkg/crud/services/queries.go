package services

const GetBurgers = "SELECT id, name, price, description FROM burgers WHERE removed = FALSE;"
const SaveBurger = "INSERT INTO burgers (name, price, description) VALUES ($1, $2, $3);"
const RemoveBurger = "UPDATE burgers SET removed = TRUE WHERE id = $1;"
