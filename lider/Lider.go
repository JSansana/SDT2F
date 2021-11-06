package main

import (
	"context"
	"log"

	"math/rand"
	"net"
	"time"

	//"github.com/streadway/amqp"

	pb "github.com/JSansana/SDT2F/proto"
	"google.golang.org/grpc"
)

const (
	puertoJL = ":50051"
)

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

type Comunicacion_Jugador_LiderServer struct {
	pb.UnimplementedComunicacion_Jugador_LiderServer
}

type Comunicacion_Lider_NamenodeServer struct {
	pb.UnimplementedComunicacion_Lider_NamenodeServer
}

/*
func MandarPozo(Arreglo []int32,address string, etapa int32){

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()

	//nuevo cliente
	c := pb.NewComunicacion_Lider_Pozo(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//r guarda la respuesta del servidor
	r, err := c.Arreglo_Pozo(ctx, &pb.ArregloLider_A_Pozo{Muertos: Arreglo, Etapa:etapa})
	if err != nil {
		log.Fatalf("No se pudo enviar solicitud: %v", err)
	}
}

func comunicacionLN(Plays []int32 , nEtapa int32){
	conn, err := grpc.Dial(puertoLN, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()

	//nuevo cliente
	c := pb.Comunicacion_Lider_Namenode(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//r guarda la respuesta del servidor
	r, err := c.Enviar_Jugada(ctx, &pb.Lider_A_Namenode{Jugadas: Plays, Etapa: nEtapa })
	if err != nil {
		log.Fatalf("No se pudo enviar solicitud: %v", err)
	}
	Respuesta := r.GetRespuesta()
	return Respuesta
}

func comunicacionLP(){
	conn, err := grpc.Dial(puertoLP, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("No conectó: %v", err)
	}
	defer conn.Close()

	//nuevo cliente
	c := pb.NewComunicacion_Lider_Pozo(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//r guarda la respuesta del servidor
	r, err := c.Pozo_Lider(ctx, &pb.Lider_A_Pozo{Solicitud: 1})
	if err != nil {
		log.Fatalf("No se pudo enviar solicitud: %v", err)
	}
	MontoAcumulado := r.GetMonto_Acumulado()
	return MontoAcumulado
}
*/
func (s *Comunicacion_Jugador_LiderServer) Etapa1_Mayor(ctx context.Context, in *pb.Jugador_A_Lider_E1_Mayor) (*pb.Lider_A_Jugador_E1_Mayor, error) {
	rand.Seed(time.Now().UnixNano())
	min := 6
	max := 10
	var Jugada_Lider int32 = int32(rand.Intn(max-min+1) + min)
	log.Printf("---------------------------------")
	log.Printf("EL lider elige %v cualquier jugador que escoja %v o más será eliminado", Jugada_Lider, Jugada_Lider)
	log.Printf("---------------------------------")
	Plays := in.GetJugadas()
	muertos := []int32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	for i := 0; i < 16; i++ {
		//AQUI SE USARA EL NAMENODE enviando Etapa 1, Plays[i], y Jugador i + 1, descartando los -1
		if Plays[i] < Jugada_Lider && Plays[i] != -1 {
			log.Printf("El jugador %v elige %v", i+1, Plays[i])
			muertos[i] = int32(i + 1)
		} else {
			if Plays[i] != -1 {
				log.Printf("El jugador %v elige %v, es eliminado", i+1, Plays[i])
			}
			muertos[i] = 0

		}
	}
	return &pb.Lider_A_Jugador_E1_Mayor{Eliminados: muertos}, nil
}

func (s *Comunicacion_Jugador_LiderServer) Cuatro(ctx context.Context, in *pb.Jugador_A_Lider_E1_21) (*pb.Lider_A_Jugador_E1_21, error) {
	log.Printf("---------------------------------")
	muertos := []int32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	Scores := in.GetPuntajes()
	for i := 0; i < 16; i++ {
		if Scores[i] < 21 {
			if Scores[i] != -1 {
				log.Printf("El jugador %v suma menos de 21, es eliminado", i+1)
			}
			muertos[i] = 0

		} else {
			log.Printf("Jugador %v consiguió más de 21 puntos", i+1)
			muertos[i] = int32(i + 1)
		}
	}
	log.Printf("---------------------------------")
	return &pb.Lider_A_Jugador_E1_21{Eliminados: muertos}, nil
}

func (s *Comunicacion_Jugador_LiderServer) Etapa2(ctx context.Context, in *pb.Jugador_A_Lider_E2_Cuerda) (*pb.Lider_A_Jugador_E2_Cuerda, error) {

	score1 := in.GetPuntaje1()
	score2 := in.GetPuntaje2()
	team1 := in.GetEquipo1()
	team2 := in.GetEquipo2()

	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 4
	var Jugada_Lider int32 = int32(rand.Intn(max-min+1) + min)
	log.Printf("---------------------------------")
	log.Printf("El lider juega: %v", Jugada_Lider)
	log.Printf("---------------------------------")
	log.Printf("Puntaje obtenido por el equipo 1: %v", score1)
	log.Printf("Puntaje obtenido por el equipo 2: %v", score2)
	log.Printf("---------------------------------")
	muertos := []int32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < len(muertos); i++ {
		muertos[i] = 0

	}
	if Jugada_Lider%2 == 0 {
		log.Printf("El numero del lider es par, por lo tanto:")
		log.Printf("---------------------------------")

		if score1%2 == 0 && score2%2 == 1 {
			for i := 0; i < len(muertos) && team1[i] != -1; i++ {
				muertos[team1[i]-1] = team1[i]
			}
			log.Printf("Equipo 2 eliminado")

		} else if score1%2 == 1 && score2%2 == 0 {
			for i := 0; i < len(muertos) && team1[i] != -1; i++ {
				muertos[team2[i]-1] = int32(team2[i])
			}
			log.Printf("Equipo 1 eliminado")

		} else if score1%2 == 0 && score2%2 == 0 {
			for i := 0; i < len(muertos) && team1[i] != -1; i++ {
				muertos[team1[i]-1] = int32(team1[i])
				muertos[team2[i]-1] = int32(team2[i])

			}
			log.Printf("Ambos equipos pasan!")

		} else if score1%2 == 1 && score2%2 == 1 {
			rand.Seed(time.Now().UnixNano())
			min := 1
			max := 2
			var opcion int32 = int32(rand.Intn(max-min+1) + min)
			//sobrevive team 1
			if opcion == 1 {
				for i := 0; i < len(muertos) && team1[i] != -1; i++ {
					muertos[team1[i]-1] = int32(team1[i])
				}

				log.Printf("Equipo 2 eliminado")
				//Sobrevive team 2
			} else if opcion == 2 {
				for i := 0; i < len(muertos) && team1[i] != -1; i++ {
					muertos[team2[i]-1] = int32(team2[i])
				}
				log.Printf("Equipo 1 eliminado")
			}
		}

	} else if Jugada_Lider%2 == 1 {
		log.Printf("El numero del lider es impar, por lo tanto:")
		log.Printf("---------------------------------")

		if score1%2 == 0 && score2%2 == 1 {
			for i := 0; i < len(muertos) && team1[i] != -1; i++ {
				muertos[team2[i]-1] = int32(team2[i])
			}
			log.Printf("Equipo 1 eliminado")

		} else if score1%2 == 1 && score2%2 == 0 {
			for i := 0; i < len(muertos) && team1[i] != -1; i++ {
				muertos[team1[i]-1] = int32(team1[i])
			}
			log.Printf("Equipo 2 eliminado")

		} else if score1%2 == 0 && score2%2 == 0 {
			for i := 0; i < len(muertos) && team1[i] != -1; i++ {
				rand.Seed(time.Now().UnixNano())
				min := 1
				max := 2
				var opcion int32 = int32(rand.Intn(max-min+1) + min)
				//sobrevive team 1
				if opcion == 1 {
					for i = 0; i < len(muertos) && team1[i] != -1; i++ {
						muertos[team1[i]-1] = int32(team1[i])

					}

					log.Printf("Equipo 2 eliminado")
					//sobrevive team 2
				} else if opcion == 2 {
					for i = 0; i < len(muertos) && team1[i] != -1; i++ {
						muertos[team2[i]-1] = int32(team2[i])

					}
					log.Printf("Equipo 1 eliminado")
				}
			}

		} else if score1%2 == 1 && score2%2 == 1 {
			for i := 0; i < len(muertos) && team1[i] != -1; i++ {
				muertos[team1[i]-1] = int32(team1[i])
				muertos[team2[i]-1] = int32(team2[i])

			}
			log.Printf("Ambos equipos pasan!")
		}
	}

	return &pb.Lider_A_Jugador_E2_Cuerda{Eliminados: muertos}, nil
}

func (s *Comunicacion_Jugador_LiderServer) Etapa3(ctx context.Context, in *pb.Jugador_A_Lider_E3) (*pb.Lider_A_Jugador_E3, error) {
	Jugadas := in.GetJugadas()
	Equipos := in.GetEquipos()
	muertos := []int32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 10
	var Jugada_Lider int32 = int32(rand.Intn(max-min+1) + min)

	var Puntaje_1 int32
	var Puntaje_2 int32

	for i := 0; i < 16 && Jugadas[2*i] != -1; i++ {
		Puntaje_1 = Abs(Jugada_Lider - Jugadas[2*i])
		Puntaje_2 = Abs(Jugada_Lider - Jugadas[2*i+1])

		//AQUI IMPLEMENTAR EL NAMENODE

		if Jugadas[2*i] == Jugadas[2*i+1] {
			log.Printf("Jugador %v y Jugador %v eligieron el mismo numero, ambos ganan", Equipos[2*i], Equipos[2*i+1])
		} else if Puntaje_1 > Puntaje_2 {
			log.Printf("Jugador %v obtiene %v y Jugador %v obtiene %v, Jugador %v gana", Equipos[2*i], Puntaje_1, Equipos[2*i+1], Puntaje_2)
			muertos[Equipos[2*i]-1] = Equipos[2*i]

		} else {
			log.Printf("Jugador %v obtiene %v y Jugador %v obtiene %v, Jugador %v gana", Equipos[2*i], Puntaje_1, Equipos[2*i+1], Puntaje_2)
			muertos[Equipos[2*i+1]-1] = Equipos[2*i+1]

		}
	}

	log.Printf("El juego del Calamar ha terminado!")

	log.Printf("Felicitaciones! Gracias por jugar, un abrazo!")

	//AQUI EL POZO

	return &pb.Lider_A_Jugador_E3{Eliminados: muertos}, nil
}

func (s *Comunicacion_Jugador_LiderServer) Fin_Etapa(ctx context.Context, in *pb.Jugador_A_Lider_Fin) (*pb.Lider_A_Jugador_Fin, error) {

	EtapaActual := in.GetEtapa_Actual()

	if EtapaActual == 0 {
		log.Printf("---------------------------------")
		log.Printf("Los jugadores pasan a la Etapa 1")
		log.Printf("---------------------------------")
		return &pb.Lider_A_Jugador_Fin{Respuesta: 1}, nil

	} else if EtapaActual == 1 {
		log.Printf("---------------------------------")
		log.Printf("Los jugadores pasan a la Etapa 2")
		log.Printf("---------------------------------")
		return &pb.Lider_A_Jugador_Fin{Respuesta: 2}, nil
	} else {
		log.Printf("---------------------------------")
		log.Printf("Los jugadores pasan a la Etapa 3")
		log.Printf("---------------------------------")
		return &pb.Lider_A_Jugador_Fin{Respuesta: 3}, nil
	}
}

func (s *Comunicacion_Jugador_LiderServer) Eliminar(ctx context.Context, in *pb.Jugador_A_Lider_Eliminar) (*pb.Lider_A_Jugador_Eliminar, error) {
	eliminado := in.GetEliminado()

	//MANDAR eliminado al POZO por RABBIT
	log.Printf("Debido al azar, el jugador %v ha sido eliminado", eliminado)
	log.Printf("---------------------------------")

	return &pb.Lider_A_Jugador_Eliminar{Respuesta: int32(0)}, nil
}

/*
func (s *Comunicacion_Jugador_LiderServer) Pozo_Jugador(ctx context.Context, in *pb.Jugador_A_Lider_Pozo) (*pb.Lider_A_Jugador_Pozo, error){

	solicitud := in.GetSolicitud()
	//PEDIR MONTO TOTAL A POZO

	log.Printf("El monto total es %v ", montopozo)
	return &pb.Lider_A_Jugador_Pozo {Monto_Acumulado : montopozo}
} */

func main() {
	lis, err := net.Listen("tcp", puertoJL)
	if err != nil {
		log.Fatalf("Fallo al escuchar: %v", err)
	}
	log.Printf("---------------------------------")
	log.Printf("Bienvenido al juego del calamar!")
	log.Printf("---------------------------------\n")
	s := grpc.NewServer()
	pb.RegisterComunicacion_Jugador_LiderServer(s, &Comunicacion_Jugador_LiderServer{})
	log.Printf("Servidor escuchando en %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo en serve: %v", err)
	}
}
