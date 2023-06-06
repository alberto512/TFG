type Account = {
  id: string
  iban: string
  name: string
  currency: string
  amount: number
  bank: string
}

type Category = {
  id: string
  name: string
}

type Transaction = {
  id: string
  description: string
  date: number
  amount: number
  category: Category
}