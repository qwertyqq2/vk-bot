## Simple VK-chat

### Установка

     go get github.com/qwertyqq2/vk-bot


### Запуск

     git clone github.com/qwertyqq2/vk-bot

     cd vk-bot

Далее измените добавьте ваш токен и айди группы в env файл - *configns/envs/conf.env*

Запуск бота

     ./bot

     

### Использование

#### Создание бота

     // Создание бота
     // GroupID - айди группы
     // Token - токен группы
     // Debug - дебаг мод
     mybot := bot.NewBot(
               bot.GroupID("groupID"),
               bot.Token("token"),
               bot.Debug(true),
          )

     // Cтарт бота
     mybot.Init()

#### Добавление колбэков

     //Приветственная кнопка
     begin := bot.NewInitCallback("text")

     //Кнопки второго уровня
	item1 := bot.NewCallback("button11", "text11")
	item2 := bot.NewCallback("button12", "text12")

     //Добавить колбэки для кнопок
     begin.AddNext(item1, item2)
     mybot.Build(begin)

#### Методы создание колбэков:

     //Создать приветственный колбэк
     NewInitCallback(text string)

     //Создать новый колбэк
     NewCallback(name, text string) 

     //Создать колбэк с выводом сообщения
     NewCallbackWithMessage(name, mes string)

     //Создать отложенный колбэк с выводом сообщения
     NewWaitCallbackWithMessage(name, mes, text string)     



### Примечание

Бот не поддерживает базу данных, поэтому, после перезапуска бота, для корректной работы бота нужно ввести приветственное сообщение пользователю.