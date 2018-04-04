package lib

import (
  "fmt"
)

func RemoveByColumn(values []string, dataset [][]string, col int) [][]string {
  for _, value := range values {
    for i, row := range dataset {
      if row[col] == value {
        dataset = append(dataset[:i], dataset[i+1:]...)

        fmt.Printf("Deleting %s\n", value) 
        break
      }
    }
  }

  return dataset
}
