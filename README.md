# Magus
Тестовое задание

Задача: реализовать универсальный метод поочередного однократного гарантированного 
выполнения входящей функции в отдельном потоке вне зависимости от того, как много 
произойдет вызовов в момент выполнения запущенной в последний раз функции.

Пример

Задачу легче понять на конкретном примере. Есть функция:

func fn() {
   	 time.Sleep(time.Second)	// Имитация 1 секунды длительности работы
   	 glog.Debug("fn")		// Вывод в лог факта завершения выполнения
}

Есть цикл длительностью 1530 миллисекунд, за каждый проход запускающий 
асинхронно или в отдельном потоке функцию fn() (не блокирующий запуск)

func main() {
     ch := time.After(time.Millisecond * 1530)
LOOP:
    for {
   	 select {
   	 case <-ch:
   		 glog.Debug("exit")
   		 break LOOP
   	 default:
   		 go fn()				// Вызов функции в отдельном потоке
   	 }
    }
    time.Sleep(time.Second * 2)
}

Задача сделать так, чтобы в один момент времени функция fn() выполнялась 
только один раз, но при этом если в момент выполнения происходят попытки 
запуска, то по завершении выполнения fn() должна запуститься еще один раз, 
вне зависимости от кол-ва обращений.


Лог будет выглядеть так:
13:49:44.489 ▶ fn
13:49:45.017 ▶ exit
13:49:45.491 ▶ fn
13:49:46.491 ▶ fn

Реализованный метод должен быть универсальным и работать с любыми функциями,
которые будут в него переданы. Выполнение функций должно происходить в 
отдельном потоке.

Задание можно выполнить на удобном для вас языке программирования.
