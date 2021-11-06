package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	pb "github.com/JSansana/SDT2F/proto"
	"google.golang.org/grpc"
)

const (
	direccionJL = "dist21:50051"
)

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func solicitudfase(fase int32, direccionJL string) int32 {

	conn, err := grpc.Dial(direccionJL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()

	//nuevo cliente
	c := pb.NewComunicacion_Jugador_LiderClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Fin_Etapa(ctx, &pb.Jugador_A_Lider_Fin{Etapa_Actual: fase})
	if err != nil {
		log.Fatalf("No se pudo enviar solicitud: %v", err)
	}

	log.Printf("Respuesta del lider : %v", r.GetRespuesta())

	respuesta := r.GetRespuesta()
	return respuesta

}

func eliminarjugador(eliminado int32, direccionJL string) int32 {
	conn, err := grpc.Dial(direccionJL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()

	c := pb.NewComunicacion_Jugador_LiderClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Eliminar(ctx, &pb.Jugador_A_Lider_Eliminar{Eliminado: eliminado})
	if err != nil {
		log.Fatalf("No se pudo enviar solicitud: %v", err)
	}
	el := r.GetRespuesta()
	return el
}

/*
func ver_pozo(Players []int32) {
	// LLamar al pozo para ver el resultado
	//ACA SE REALIZA JUGADOR_LIDER para solicitar el monto del pozo.

	conn, err := grpc.Dial(direccionJL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()

	//nuevo cliente
	c := pb.NewComunicacion_Jugador_LiderClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Pozo_Jugador(ctx, &pb.Jugador_A_Lider_Pozo{Solicitud: 2})
	if err != nil {
		log.Fatalf("No se pudo enviar solicitud: %v", err)
	}

	monto := r.GetMonto_Acumulado()
	fmt.Println("El pozo actual es de: %v ",)
	return 0
}
*/
func Etapa3(Players []int32) int32 {
	var inputint int32

	for {
		//Pregunta al jugador si quiere entrar al juego o si antes quiere ver el pozo.
		log.Printf("---------------------------------")
		fmt.Println("¿Qué quieres hacer? Ingresa el numero de la opcion para continuar")
		fmt.Println("1) Jugar al Todo o Nada")

		fmt.Scanln(&inputint)

		/*
			if inputint == 3 {
				//Enviar solicitud de revision de jugadas, recibe

			}
		*/

		if inputint == 2 {
			//ver_pozo()
			break
		}
		if inputint == 1 {
			break
			//Avisar al lider de que ya estamos jugando
		}
	}
	log.Printf("---------------------------------")
	//Conteo de jugadores vivos
	count := 0
	var Jugador int32 = 1
	for i := 0; i < 16; i++ {
		if Players[i] != 0 {
			count++
			if Players[i] != 1 {
				Jugador = Players[i]
			}
		}
	}

	if count == 2 {
		log.Printf("---------------------------------")

		fmt.Println("Escoge un número, del 1 al 10.")
		for {
			fmt.Scanln(&inputint)
			if inputint >= 1 && inputint <= 10 {
				break
			} else {
				fmt.Println("Número fuera de rango, por favor del 1 al 10!")
			}
		}

		rand.Seed(time.Now().UnixNano())
		min := 0
		max := 10
		var OtroJugador int32 = int32(rand.Intn(max-min+1) + min)
		var Lider int32 = int32(rand.Intn(max-min+1) + min)

		log.Printf("---------------------------------")
		log.Printf("El lider juega %v", Lider)
		log.Printf("---------------------------------")

		puntaje1 := Abs(Lider - OtroJugador)
		puntaje2 := Abs(Lider - inputint)

		log.Printf("El jugador %v consigue %v", Jugador, puntaje1)
		log.Printf("El jugador 1 consigue %v", puntaje2)

		if puntaje1 == puntaje2 {
			log.Printf("Ambos ganan, felicidades!")
		} else if puntaje1 > puntaje2 {
			log.Printf("El jugador %v gana, has sido eliminado", Jugador)
		} else {
			log.Printf("Has ganado el juego del calamar!")
		}

		return 0
	}

	//Si la cantidad de jugadores es impar, se mata uno al azar
	if count%2 != 0 {
		for {
			rand.Seed(time.Now().UnixNano())
			min := 0
			max := 15
			var randomint int32 = int32(rand.Intn(max-min+1) + min)
			if Players[randomint] != 0 {
				Players[randomint] = 0
				// Jugador a Lider : enviar randomint a lider para que elimine a ese jugador
				res := eliminarjugador(randomint, direccionJL)
				log.Printf("Lider responde %v", res)
				count--
				break
			}
		}
	}

	//Si el jugador escogido somos nosotros, termina el juego
	if Players[0] == 0 {
		log.Printf("Debido al azar, has muerto")
		return 0

	}

	//Se crea una lista con los jugadores disponibles y un arreglo "de arreglos con las parejas"
	/*
		Availables := make([]int32, count)
		Teams := make([][]int32, count/2, 2) */
	Availables := [16]int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
	Teams := [8][2]int32{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}

	//Se lista a los jugadores disponibles
	j := 0
	for i := 0; i < 16; i++ {
		if Players[i] != 0 {
			Availables[j] = Players[i]
			j++
		}
	}

	//Se asignan aleatoriamente las parejas
	k := 0
	for i := 0; i < count/2; i++ {
		for {
			rand.Seed(time.Now().UnixNano())
			min := 0
			max := count - 1
			var randomint int32 = int32(rand.Intn(max-min+1) + min)
			//Se confirma si el jugador no fue asignado antes
			if Availables[randomint] != -1 {
				Teams[i][j] = Availables[randomint]
				Availables[randomint] = 0
				//A medida que se asignan los jugadores, avanza el arreglo
				if k == 0 {
					k++
				} else {
					k--
					i++
				}
				break
			}
		}
	}

	//Crea un arreglo con los puntajes de cada pareja
	Scores := make([][]int32, count/2, 2)
	log.Printf("---------------------------------")

	fmt.Println("Escoge un número, del 1 al 10.")
	for {
		fmt.Scanln(&inputint)
		if inputint >= 1 && inputint <= 10 {
			break
		} else {
			fmt.Println("Número fuera de rango, por favor del 1 al 10!")
		}
	}

	for i := 0; i < count/2; i++ {
		for j := 0; j < 2; j++ {
			//Si es el jugador se almacena el valor,
			if Teams[i][j] == 1 {
				Scores[i][j] = inputint
			} else {
				rand.Seed(time.Now().UnixNano())
				min := 0
				max := 10
				var randomint int32 = int32(rand.Intn(max-min+1) + min)
				Scores[i][j] = randomint
			}
		}
	}
	log.Printf("---------------------------------")
	jugadas := []int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}
	equipos := []int32{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1, -1}

	for i := 0; i < count/2; i++ {
		jugadas[2*i] = Scores[i][0]
		jugadas[2*i+1] = Scores[i][1]

		equipos[2*i] = Teams[i][0]
		equipos[2*i+1] = Teams[i][1]
	}

	//Jugador a lider : Enviar los dos arreglos 1D y retorna un arreglo de todos los que mueren
	conn, err := grpc.Dial(direccionJL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()

	c := pb.NewComunicacion_Jugador_LiderClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Etapa3(ctx, &pb.Jugador_A_Lider_E3{Jugadas: jugadas, Equipos: equipos})
	if err != nil {
		log.Fatalf("No se pudo enviar solicitud: %v", err)
	}
	muertos := r.GetEliminados()
	for i := 0; i < len(muertos)/2; i++ {
		Players[muertos[i]-1] = 0
	}
	if Players[0] != 0 {
		fmt.Println("Has ganado, felicidades.")
		//Mostrar pozo ganador.
	}
	return 0
}

func Etapa2(Players []int32) int32 {

	var inputint int32
	for {
		//Pregunta al jugador si quiere entrar al juego o si antes quiere ver el pozo.
		log.Printf("---------------------------------")
		fmt.Println("¿Qué quieres hacer? Ingresa el numero de la opcion para continuar")
		fmt.Println("1) Jugar a Tirar la Cuerda")

		fmt.Scanln(&inputint)

		if inputint == 2 {
			//ver_pozo()
		}
		if inputint == 1 {
			break
		}
	}

	//Conteo de jugadores vivos
	count := 0
	for i := 0; i < 16; i++ {
		if Players[i] != 0 {
			count++
		}
	}
	log.Printf("---------------------------------")

	//Si la cantidad de jugadores es impar, se mata uno al azar
	if count%2 != 0 {
		for {
			rand.Seed(time.Now().UnixNano())
			min := 0
			max := 15
			var randomint int32 = int32(rand.Intn(max-min+1) + min)
			if Players[randomint] != 0 {
				Players[randomint] = 0
				// Jugador a Lider : enviar randomint a lider para que elimine a ese jugador
				eliminarjugador(randomint+1, direccionJL)
				count--
				break
			}
		}
	}

	//Si el jugador escogido somos nosotros, termina el juego
	if Players[0] == 0 {
		log.Printf("Has muerto!")
		return 0
	}

	//Se hace una lista de jugadores disponibles y otras dos con los equipos
	Availables := make([]int32, count)
	Team1 := []int32{-1, -1, -1, -1, -1, -1, -1, -1}
	Team2 := []int32{-1, -1, -1, -1, -1, -1, -1, -1}

	//Se lista a los jugadores disponibles
	j := 0
	for i := 0; i < 16; i++ {
		if Players[i] != 0 {
			Availables[j] = Players[i]
			j++
		}
	}

	log.Printf("---------------------------------")
	log.Printf("El equipo 1 se compone de:")
	log.Printf("---------------------------------")
	//Se asigna aleatoriamente jugadores al equipo 1.
	for i := 0; i < count/2; i++ {
		for {
			rand.Seed(time.Now().UnixNano())
			min := 0
			max := count - 1
			var randomint int32 = int32(rand.Intn(max-min+1) + min)
			//Revisa si el jugador no fue asignado antes
			if Availables[randomint] != 0 {
				log.Printf("Jugador %v", Availables[randomint])
				Team1[i] = Availables[randomint]
				Availables[randomint] = 0
				break
			}
		}
	}

	log.Printf("---------------------------------")
	log.Printf("El equipo 2 se compone de:")
	log.Printf("---------------------------------")
	//Se asigna el resto de jugadores al equipo 2
	contador := 0
	for i := 0; i < count; i++ {
		if Availables[i] != 0 {
			Team2[contador] = Availables[i]
			log.Printf("Jugador %v", Availables[i])
			contador++
		}
	}

	//Puntajes de cada equipo
	var Scores1 int32 = 0
	var Scores2 int32 = 0

	//El jugador escoge un numero
	log.Printf("---------------------------------")
	fmt.Println("Escoge un número, del 1 al 4.")
	//Se confirma si el numero esta dentro del rango permitido
	for {
		if inputint >= 1 && inputint <= 4 {
			fmt.Scanln(&inputint)
			//Revisa a qué equipo pertenece el jugador y asigna su respectivo puntaje
			for j := 0; j < count/2; j++ {
				if 1 == Team1[j] {
					Scores1 += inputint
					break
				}
				if 1 == Team2[j] {
					Scores2 += inputint
					break
				}
			}
			break
		} else {
			fmt.Println("Número fuera de rango, por favor del 1 al 4!")
		}
	}
	log.Printf("---------------------------------")
	for i := 0; i < count/2; i++ {
		rand.Seed(time.Now().UnixNano())
		min := 1
		max := 4
		var randomint int32 = int32(rand.Intn(max-min+1) + min)

		if Team1[i] != 1 {
			Scores1 += randomint
		}
		if Team2[i] != 1 {
			Scores2 += randomint
		}

	}

	//Jugador a Lider : Enviar los dos puntajes obtenidos por cada equipo, además de ambos equipos, retorna un arreglo con los muertos.
	//muertos[]
	//
	conn, err := grpc.Dial(direccionJL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()

	c := pb.NewComunicacion_Jugador_LiderClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Etapa2(ctx, &pb.Jugador_A_Lider_E2_Cuerda{Puntaje1: Scores1, Puntaje2: Scores2, Equipo1: Team1, Equipo2: Team2})
	if err != nil {
		log.Fatalf("No se pudo enviar solicitud: %v", err)
	}
	Players = r.GetEliminados()

	if Players[0] == 0 {
		log.Printf("Has muerto!")
		return 0
	}

	/*
		for i := 0; i < len(muertos)/2; i++ {
			if muertos[i] != 0 {
				Players[muertos[i]-1] = 0
			}

		} */

	//Conteo de vivos
	count2 := 0
	for i := 0; i < 16; i++ {
		if Players[i] != 0 {
			count2++
		}
	}

	//Si todos los jugadores mueren, se acaba el juego
	if count2 == 0 {
		fmt.Println("Todos los jugadores han muerto! se termina el juego.")

	} else if count2 == 1 {
		fmt.Println("Solo tu quedas vivo, felicidades, has ganado!")
	} else if count2 == 3 {
		fmt.Println("Quedan 3 jugadores vivos, por lo tanto no podrán jugar a la Etapa 3, todos ganan!")

	} else {
		respuesta := solicitudfase(2, direccionJL)
		if respuesta == 3 {
			Etapa3(Players)
		} else {
			fmt.Println("Adiós.")
		}
	}
	return 0
}

func Etapa1(Players []int32) int32 {

	//
	var inputint int32

	//Pregunta al jugador si quiere entrar al juego o si quiere antes ver el pozo
	for {
		log.Printf("---------------------------------")
		fmt.Println("¿Qué quieres hacer? Ingresa el numero de la opcion para continuar")
		fmt.Println("1) Jugar Luz Roja, Luz Verde")

		fmt.Scanln(&inputint)

		if inputint == 2 {
			//ver_pozo()

		}
		if inputint == 1 {
			break
		}
	}

	//Crea un arreglo con los puntajes de cada jugador
	Scores := []int32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	//Arreglo con la jugada de dicha ronda
	Plays := []int32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	//For para las rondas
	for rounds := 1; rounds <= 4; rounds++ {

		//Pregunta por el puntaje
		log.Printf("---------------------------------")
		fmt.Println("Escoge un número, del 1 al 10.")

		//Se confirma si el numero escogido está dentro del rango permitido

		for {
			fmt.Scanln(&inputint)
			if inputint >= 1 && inputint <= 10 {
				Scores[0] += inputint
				Plays[0] = inputint
				break
			} else {
				fmt.Println("Número fuera de rango, por favor del 1 al 10!")
			}
		}

		log.Printf("---------------------------------")
		//Puntaje de los bots
		for i := 1; i < 16; i++ {
			if Players[i] != 0 {
				rand.Seed(time.Now().UnixNano())
				min := 1
				max := 10
				var randomint int32 = int32(rand.Intn(max-min+1) + min)
				Scores[i] += randomint
				Plays[i] = randomint
				log.Printf("Jugador: %v escoge %v", Players[i], randomint)
				time.Sleep(100 * time.Millisecond)
			}
		}

		//Jugador a lider : Se envia los arreglos con las jugadas de la ronda, y el lider responde con los jugadores eliminados
		// Lider devuelve un arreglo con los 16 jugadores y un 0 en la posicion en que un jugador haya muerto.

		conn, err := grpc.Dial(direccionJL, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("No conectó: %v", err)
		}
		defer conn.Close()

		c := pb.NewComunicacion_Jugador_LiderClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.Etapa1_Mayor(ctx, &pb.Jugador_A_Lider_E1_Mayor{Jugadas: Plays})
		if err != nil {
			log.Fatalf("No se pudo enviar solicitud: %v", err)
		}
		Players = r.GetEliminados()
		/*
			for i := 1; i < 16; i++ {
				if arreglo_respuesta[i] == 0 {
					Players[i] = 0 // El arreglo original de los jugadores es modificado con la respuesta del Lider, es decir, se matan a los jugadores
					//Se Utiliza un -1 para que esas posiciones del arreglo no influyan en las jugadas futuras
					Plays[i] = -1
					Scores[i] = -1
				}
			}*/
		if Players[0] == 0 {
			log.Printf("Has muerto porque tu numero fue mayor o igual que el del lider!")
			return 0

		}

	}
	//Jugador a lider : Enviar arreglo de puntaje total y el lider devuelve un arreglo con los jugadores que hayan muerto

	conn, err := grpc.Dial(direccionJL, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()

	c := pb.NewComunicacion_Jugador_LiderClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Cuatro(ctx, &pb.Jugador_A_Lider_E1_21{Puntajes: Scores})
	if err != nil {
		log.Fatalf("No se pudo enviar solicitud: %v", err)
	}
	arreglo_respuesta := r.GetEliminados()
	//Conteo de vidas y actualizacion de array players. Si todos murieron, se detiene el juego
	count := 0
	for i := 0; i < 16; i++ {
		if arreglo_respuesta[i] == 0 {
			Players[i] = 0
		} else {
			count++

		}
	}
	if count == 0 {
		fmt.Println("Todos los jugadores murieron")
		return 0
	} else if count == 1 {
		fmt.Println("Terminó el juego, solo sobrevivió 1 jugador")
		return 0

	} else {
		//Jugador a Lider : Solicitud de entrada a etapa 2, el lider indica el inicio de la siguiente etapa
		respuesta := solicitudfase(1, direccionJL)

		if respuesta == 2 {
			Etapa2(Players)
			return 0
		} else {
			fmt.Println("Adiós.")
			return 0
		}

	}
}
func main() {

	//var inputstr string

	//Confirmacion a la entreada del juego
	log.Printf("---------------------------------")
	fmt.Println("¿Quieres entrar al juego del calamar?")
	fmt.Println("Ingrese 'si' para solicitar entrar al juego")
	var inputstr string
	fmt.Scanln(&inputstr)

	//Si la respuesta es positiva, entra el juego.
	if inputstr == "Si" || inputstr == "si" {

		Players := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

		//Aqui se solicitad al Lider entrar al juego del calamar
		conn, err := grpc.Dial(direccionJL, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("No conectó: %v", err)
		}
		defer conn.Close()

		//nuevo cliente
		c := pb.NewComunicacion_Jugador_LiderClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// ACA SE REALIZA JUGADOR_A_LIDER para solicitar entrar al juego
		//r guarda la respuesta del servidor
		r, err := c.Fin_Etapa(ctx, &pb.Jugador_A_Lider_Fin{Etapa_Actual: 0})
		if err != nil {
			log.Fatalf("No se pudo enviar solicitud: %v", err)
		}

		log.Printf("El lider permite tu entrada! %v", r.GetRespuesta())
		respuesta := r.GetRespuesta()
		if respuesta == 1 {
			log.Printf("---------------------------------")
			fmt.Println("Bienvenido entonces, al juego del calamar.\n")
			Etapa1(Players)

		}

	} else {
		fmt.Println("Adiós.")
	}
	return
}
