package main

import(
  "math/rand"
  "fmt"
  "time"
  "math"
  "strconv"
  "strings"

)

const(
  numberByte = 6 // 6, 3
  lenSub = 64 // 64, 8
  getBitlen = 5 // 5, 2

  xor = `⊕`
)

func main(){
  start := time.Now()

  vector := gen()
  fmt.Println(vector)
  f := make([][]byte, numberByte)

  j := 0
  for i := getBitlen; i >= 0; i-- {
    f[j] = getCoordVec(vector, i)
    j++
  }

  weightBF := make([]int, numberByte)
  for i := 0; i < numberByte; i++{
    fmt.Printf("f[%d] = ",i+1)
    fmt.Println(f[i])
    weightBF[i] = boolFuncWeight(f[i])
  }

  fmt.Println("||f|| = ", weightBF)

  for i := range f{
    polBF, polStr := fft(f[i])
    fmt.Printf("Полином жегалкина для f[%d] = ", i+1)
    fmt.Println(polBF)
    /*var b = make([]byte, len(polStr)-3) ?? //удаление ? в конце
    copy(b, polStr)??
    fmt.Println(string(b))*/
    fmt.Println("с xor",polStr)
    fmt.Println("без xor", polStr[:len(polStr)-3])
    fict := searchFictVar(polStr)

    switch fict {
    case "":
      fmt.Println("Фиктивные переменные отсутствуют")
    default:
      fmt.Println(fict, " - фиктивные переменные")
    }
  }

  duration := time.Since(start)
  fmt.Println(duration)
  fmt.Println(duration.Nanoseconds(), " - наносекунд")
}
//вес бф f_i
func boolFuncWeight(bf []byte) int{
  weight := 0
  for i := range bf{
    if bf[i] == 1{
      weight++
    }
  }
  return weight
}
// вектор значения координатной функции отображения S
func getCoordVec(vec []byte, i int)[]byte{
  f_i := make([]byte, lenSub)
  for j := range f_i{
    f_i[j] = vec[j] >> i & 1
  }
  return f_i
}
// генерация подставновки
func gen()[]byte{
  generator := rand.New(rand.NewSource(time.Now().UnixNano()))
  vec := make([]byte, 0, lenSub)

  for {
    temp := byte(generator.Intn(lenSub))
    if add := check(vec, temp); add != false{
      vec = append(vec, temp)
    }
    if len(vec) == lenSub{
      break
    }
  }
  return vec
}
//проверяем, есть ли число в срезе
func check(vec []byte, temp byte) bool{
  for i := 0; i < len(vec); i++{
    if vec[i] == temp{
      return false
    }
  }
  return true
}
// бпф
func fft(bf []byte) ([]byte, string){
  maxLenBlock := lenSub / 2 // максимальный размер блока
  numberParts := len(bf) // количество частей в разбении алгоритма

  var j float64 = 0
  for i := 0; i <= maxLenBlock; i=int(math.Pow(2, j)){
    bf = stepFFT(bf, i, numberParts)
    numberParts = numberParts / 2 // количество блоков
    j++
  }
  return bf, getPolinom(bf)
}
// шаг алгоритма
func stepFFT(bf []byte, sizeBlock, numberParts int) []byte{
  if sizeBlock == 0{
    sizeBlock = 1
  }

  split := splitBF(bf, sizeBlock, numberParts) // делим бф на части 1,2,4 итд

  newBF := make([][]byte, numberParts)
  for i := 0; i < numberParts; i++{ // проходим по частям
    if i < numberParts{
      newBF[i] = append(newBF[i], split[i]...) // левая часть остается
      newBF[i+1] = append(newBF[i+1], XOR(split[i], split[i+1])...) // правая ксорится с левой
      i++ // переходим к след части
    }
  }

  newbf := make([]byte, 0, len(bf)) // преобразуем слайс слайсов в слайс
  for i := range newBF{
    newbf = append(newbf, newBF[i]...)
  }
  return newbf
}

func XOR(a, b []byte) []byte{
  for i := 0; i < len(a); i++{
    a[i] = a[i] ^ b[i]
  }
  return a
}

func splitBF(bf []byte, sizeBlock, numberParts int) [][]byte{
  split := make([][]byte, numberParts)

  j := 0
  for i := range split{ // разбиваем бф на части
    split[i] = append(split[i], bf[j:j+sizeBlock]...)
    j += sizeBlock
  }
  return split
}

func getPolinom(bf []byte) string{
  var polinom string
  for i := range bf{
    if bf[i] == 1 {
      monom := getMonom(i)
      polinom += monom + xor
    }
  }
  return polinom
}
// проходимя по 2 представлению позиции и возвращаем моном мн-на жегалкина
func getMonom(number int) string{ // переводим смотрим на каких позициях еденицы, собираем строку
  binaryNumber := toBinary(number)

  var monom string
  for i := range binaryNumber{
    if binaryNumber[i] == "1"{
      monom += getX(i)
    }
  }

  if monom == ""{
    monom = "1"
  }
  return monom
}
// переводим номер позиции где мн-н принимает "1" в 2сс
func toBinary(val int) []string{
  splitBinary := make([]string, 0, 6)
  for i := getBitlen; i >= 0; i--{
    splitBinary = append(splitBinary, strconv.Itoa(val >> i & 1))
  }
  return splitBinary
}
func getX(i int) string{
  parts := []string{"x1", "x2", "x3", "x4", "x5", "x6"}
  return parts[i]
}
// поиск фиктивной переменной
func searchFictVar(pZ string) string{
  i := 0
  var fict string
  for i < numberByte {
    contain := strings.Contains(pZ, getX(i))
    if !contain{
      fict += getX(i)
    }
    i++
  }
  return fict
}
