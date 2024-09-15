package main

import (
	"database/sql"
	"fmt"
	"main/foreground"
	"main/keyboard"
	"main/mouse"
	"time"

	"github.com/gonutz/w32/v2"
	_ "github.com/mattn/go-sqlite3"
)

type Resp struct {
	Id         int
	App        string
	Start_time string
	End_time   string
}

// / модуль для загрузки в caldav ( который проеобразует данные из базы данных => если события короче 15 минут он их собирает  вместе  )
// добавить чтобы при включении и выключении программы закрывались прошлые события ( + правильно если время запуска програмы раньше чем включение компьютера то поставить время полседнего выключения )
// проверка  на запуск ( чтобы не дать запустить 2 экземпляр приложения )
func closer_d(d *sql.DB) {
	// закрыть прошлую запись ( если уже не закрыта )
	var result Resp
	d.QueryRow("SELECT id,app,IFNULL(start_time,'empty'),IFNULL(end_time,'empty') FROM main ORDER BY id DESC LIMIT 1 ").Scan(&result.Id, &result.App, &result.Start_time, &result.End_time)
	if result.Id != 0 && result.End_time == "empty" {
		// fmt.Println("Не закрыт запрос предыдущий ", result)
		d.Exec("UPDATE main SET end_time = $1 WHERE id = (SELECT MAX(id) FROM main) ", time.Now().String())
	}
}

func main() {
	var app string
	unactive_cycles := 0
	db, _ := sql.Open("sqlite3", "main.db")
	db.Exec("CREATE TABLE IF NOT EXISTS 'main' (	'id'	INTEGER,	'app'	TEXT,'start_time'	TEXT,	'end_time'	TEXT,	PRIMARY KEY('id' AUTOINCREMENT))")
	keyboard.Init()
	mouse_cords_x, mouse_cords_y := 0, 0
	keyboard_key := ""
	foreground_app := ""
	fmt.Println("appclmonl starting")
	if true { // скрыть консоль
		fmt.Println("Окно скроется автоматически")
		time.Sleep(3 * time.Second)

		console := w32.GetConsoleWindow()
		if console != 0 {
			_, consoleProcID := w32.GetWindowThreadProcessId(console)
			if w32.GetCurrentProcessId() == consoleProcID {
				w32.ShowWindowAsync(console, w32.SW_HIDE)
			}
		}
	}
	fmt.Println("запуск модуля загрузки в caldav ")

	for i := 0; i < 10000000; i++ {

		local_app := foreground.Run()
		if foreground_app != local_app {
			foreground_app = local_app

			// fmt.Println("приложение изменилось")
			// закрыть предыдущий и закинуть новый
			closer_d(db)
			// !! если предыдущая запись не равна нашей
			db.QueryRow("SELECT app FROM main ORDER BY id DESC LIMIT 1 ").Scan(&app)
			if app != foreground_app {
				db.Exec("insert INTO main (app,start_time)  VALUES ($1,$2)  ", foreground_app, time.Now().String())
				unactive_cycles = 0
			}

		} else {
			local_x, local_y := mouse.Run()
			local_key := keyboard.Lastpressendkey()
			if !(mouse_cords_x == local_x && mouse_cords_y == local_y && keyboard_key == local_key) {
				mouse_cords_x = local_x
				mouse_cords_y = local_y
				keyboard_key = local_key
				//fmt.Println("Мышь сдвинулась или клавиатура нажата ")

				// unactive
				unactive_cycles = 0
				db.QueryRow("SELECT app FROM main ORDER BY id DESC LIMIT 1 ").Scan(&app)
				if app == "Not-active" {
					closer_d(db)
					db.Exec("insert INTO main (app,start_time)  VALUES ($1,$2)  ", foreground_app, time.Now().String())
				}

			} else {
				if unactive_cycles >= 10 {
					// fmt.Println("Неактивно началась")
					// закрыть прошлую и создать новую
					closer_d(db)
					db.QueryRow("SELECT app FROM main ORDER BY id DESC LIMIT 1 ").Scan(&app)
					if app != "Not-active" {
						db.Exec("insert INTO main (app,start_time)   VALUES ($1,$2)  ", "Not-active", time.Now().String())
					}
				} else {
					unactive_cycles++
				}
			}
		}

		time.Sleep(2 * time.Second)
	}

}
