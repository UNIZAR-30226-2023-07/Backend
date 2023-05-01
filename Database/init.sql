CREATE TABLE JUGADORES (
	nombre 		text NOT NULL,
    contra 		text NOT NULL,
	foto 		int NOT NULL,
    descrp 		text,
    pjugadas 	integer NOT NULL,
    pganadas 	integer NOT NULL,
	puntos 		integer NOT NULL,
    email 		text NOT NULL,
	codigo 		text,
	UNIQUE (email),
    PRIMARY KEY (codigo)
);

CREATE TABLE AMISTAD (
	id 			serial,
    estado 		text NOT NULL CHECK(estado = 'esp_confirmacion' OR estado = 'confirmada' OR estado = 'pendiente'),
    usr1 		text REFERENCES JUGADORES (codigo),
    usr2 		text REFERENCES JUGADORES (codigo),
	UNIQUE(usr1, usr2),
    PRIMARY KEY (id)
);

CREATE TABLE PARTIDAS (
	clave 		text,
	creador 	text REFERENCES JUGADORES (codigo),
	tipo 		text NOT NULL CHECK(tipo = 'amistosa' OR tipo = 'torneo'),
	estado 		text NOT NULL CHECK(estado = 'terminada' OR estado = 'pausada' OR estado = 'iniciada' OR estado = 'creando'),
	pactual 	text,
	PRIMARY KEY (clave)
);

CREATE TABLE PARTICIPAR (
	id 			serial,
	partida 	text REFERENCES PARTIDAS (clave),
	jugador 	text REFERENCES JUGADORES (codigo),
	puntos_resultado integer NOT NULL,
	enlobby 	integer NOT NULL CHECK(enlobby = 1 OR enlobby = 0),
	turno 		integer NOT NULL,
	abierto		text NOT NULL CHECK(abierto = 'si' OR abierto = 'no'),
	UNIQUE (partida, turno),
	UNIQUE (partida, jugador),
	PRIMARY KEY (id)
);

CREATE TABLE MENSAJES (
	id 			serial, 
	jug_emi 	text REFERENCES JUGADORES (codigo),
	jug_rcp 	text REFERENCES JUGADORES (codigo),
	timestamp 	timestamp NOT NULL,
	contenido 	text NOT NULL,
	leido 		integer NOT NULL,
	PRIMARY KEY(id)
);

CREATE TABLE CARTAS(
    id 			integer,
	valor 		integer NOT NULL,
    palo 		integer NOT NULL,
	reverso 	integer NOT NULL, 
	UNIQUE (valor, palo, reverso),
    PRIMARY KEY (id)
);	

CREATE TABLE COMBINACIONES (
	id 			serial,
	partida 	text REFERENCES PARTIDAS (clave),
	carta 		serial REFERENCES CARTAS (id),
	ncomb		integer NOT NULL,
	UNIQUE(partida, carta),
	PRIMARY KEY (id)
);

CREATE TABLE DESCARTES (
	id 			serial,
	partida 	text REFERENCES PARTIDAS (clave),
	carta 		serial REFERENCES CARTAS (id),
	UNIQUE (partida),
	PRIMARY KEY (id)
);

CREATE TABLE MANOS (
	id 			serial,
	partida 	text REFERENCES PARTIDAS (clave),
	carta 		serial REFERENCES CARTAS (id),
	turno 		integer NOT NULL,
	UNIQUE (partida, turno, carta),
	PRIMARY KEY (id)
);

CREATE TABLE MAZOS (
	id 			serial,
	partida 	text REFERENCES PARTIDAS (clave),
	carta 		serial REFERENCES CARTAS (id),
	UNIQUE(partida, carta),
	PRIMARY KEY (id)
);

	
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, puntos, email, codigo) values ('Max', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 7, 'Aliquam sit amet diam in magna bibendum imperdiet. Nullam orci pede, venenatis non, sodales sed, tincidunt eu, felis. Fusce posuere felis sed lacus.', 9, 14, 1000, 'max@gmail.com', 1);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, puntos, email, codigo) values ('Eugenio', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 1, 'Proin at turpis a pede posuere nonummy. Integer non velit.', 12, 87, 900, 'eu@gmail.com', 2);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, puntos, email, codigo) values ('Lauren', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 2, 'Morbi non lectus.', 94, 88, 800, 'lclipson2@trellian.com', 3);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, puntos, email, codigo) values ('Shannon', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 9, 'Cras non velit nec nisi vulputate nonummy. Maecenas tincidunt lacus at velit.', 85, 32, 700, 'shan@gmail.com', 4);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, puntos, email, codigo) values ('Valentino', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 9, 'Donec vitae nisi. Nam ultrices, libero non mattis pulvinar, nulla,', 76, 88, 600, 'val@gmail.com', 5);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, puntos, email, codigo) values ('Ernest', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 2, 'Donec posuere metus vitae ipsum. Aliquam non mauris. Morbi non lectus.', 52, 90, 500, 'erni@gmail.com', 6);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, puntos, email, codigo) values ('Frederic', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 8, 'Fusce posuere felis sed lacus. Morbi sem mauris, laoreet ut, rhoncus aliquet, pulvinar sed, nisl.', 25, 31, 400, 'fredy@gmail.com', 7);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, puntos, email, codigo) values ('Edeline', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 10, 'In hac habitasse platea dictumst. Morbi vestibulum, velit id pretium iaculis, diam erat fermentum justo, nec condimentum neque sapien placerat ante.', 95, 75, 300, 'edy@gmail.com', 8);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, puntos, email, codigo) values ('Niki', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 3, 'Nullam varius.', 29, 23, 200, 'minaj@gmail.com', 9);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, puntos, email, codigo) values ('Luce', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 2, 'Fusce consequat. Nulla nisl. Nunc nisl.', 1, 99, 100, 'luce@gamil.com', 10);

insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 2);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 3);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 4);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 5);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 6);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 7);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 8);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 9);



INSERT INTO CARTAS VALUES (011,0,1,1);
INSERT INTO CARTAS VALUES (012,0,1,2);
INSERT INTO CARTAS VALUES (021,0,2,1);
INSERT INTO CARTAS VALUES (022,0,2,2);
INSERT INTO CARTAS VALUES (031,0,3,1);
INSERT INTO CARTAS VALUES (032,0,3,2);
INSERT INTO CARTAS VALUES (041,0,4,1);
INSERT INTO CARTAS VALUES (042,0,4,2);
INSERT INTO CARTAS VALUES (111,1,1,1);
INSERT INTO CARTAS VALUES (112,1,1,2);
INSERT INTO CARTAS VALUES (121,1,2,1);
INSERT INTO CARTAS VALUES (122,1,2,2);
INSERT INTO CARTAS VALUES (131,1,3,1);
INSERT INTO CARTAS VALUES (132,1,3,2);
INSERT INTO CARTAS VALUES (141,1,4,1);
INSERT INTO CARTAS VALUES (142,1,4,2);
INSERT INTO CARTAS VALUES (211,2,1,1);
INSERT INTO CARTAS VALUES (212,2,1,2);
INSERT INTO CARTAS VALUES (221,2,2,1);
INSERT INTO CARTAS VALUES (222,2,2,2);
INSERT INTO CARTAS VALUES (231,2,3,1);
INSERT INTO CARTAS VALUES (232,2,3,2);
INSERT INTO CARTAS VALUES (241,2,4,1);
INSERT INTO CARTAS VALUES (242,2,4,2);
INSERT INTO CARTAS VALUES (311,3,1,1);
INSERT INTO CARTAS VALUES (312,3,1,2);
INSERT INTO CARTAS VALUES (321,3,2,1);
INSERT INTO CARTAS VALUES (322,3,2,2);
INSERT INTO CARTAS VALUES (331,3,3,1);
INSERT INTO CARTAS VALUES (332,3,3,2);
INSERT INTO CARTAS VALUES (341,3,4,1);
INSERT INTO CARTAS VALUES (342,3,4,2);
INSERT INTO CARTAS VALUES (411,4,1,1);
INSERT INTO CARTAS VALUES (412,4,1,2);
INSERT INTO CARTAS VALUES (421,4,2,1);
INSERT INTO CARTAS VALUES (422,4,2,2);
INSERT INTO CARTAS VALUES (431,4,3,1);
INSERT INTO CARTAS VALUES (432,4,3,2);
INSERT INTO CARTAS VALUES (441,4,4,1);
INSERT INTO CARTAS VALUES (442,4,4,2);
INSERT INTO CARTAS VALUES (511,5,1,1);
INSERT INTO CARTAS VALUES (512,5,1,2);
INSERT INTO CARTAS VALUES (521,5,2,1);
INSERT INTO CARTAS VALUES (522,5,2,2);
INSERT INTO CARTAS VALUES (531,5,3,1);
INSERT INTO CARTAS VALUES (532,5,3,2);
INSERT INTO CARTAS VALUES (541,5,4,1);
INSERT INTO CARTAS VALUES (542,5,4,2);
INSERT INTO CARTAS VALUES (611,6,1,1);
INSERT INTO CARTAS VALUES (612,6,1,2);
INSERT INTO CARTAS VALUES (621,6,2,1);
INSERT INTO CARTAS VALUES (622,6,2,2);
INSERT INTO CARTAS VALUES (631,6,3,1);
INSERT INTO CARTAS VALUES (632,6,3,2);
INSERT INTO CARTAS VALUES (641,6,4,1);
INSERT INTO CARTAS VALUES (642,6,4,2);
INSERT INTO CARTAS VALUES (711,7,1,1);
INSERT INTO CARTAS VALUES (712,7,1,2);
INSERT INTO CARTAS VALUES (721,7,2,1);
INSERT INTO CARTAS VALUES (722,7,2,2);
INSERT INTO CARTAS VALUES (731,7,3,1);
INSERT INTO CARTAS VALUES (732,7,3,2);
INSERT INTO CARTAS VALUES (741,7,4,1);
INSERT INTO CARTAS VALUES (742,7,4,2);
INSERT INTO CARTAS VALUES (811,8,1,1);
INSERT INTO CARTAS VALUES (812,8,1,2);
INSERT INTO CARTAS VALUES (821,8,2,1);
INSERT INTO CARTAS VALUES (822,8,2,2);
INSERT INTO CARTAS VALUES (831,8,3,1);
INSERT INTO CARTAS VALUES (832,8,3,2);
INSERT INTO CARTAS VALUES (841,8,4,1);
INSERT INTO CARTAS VALUES (842,8,4,2);
INSERT INTO CARTAS VALUES (911,9,1,1);
INSERT INTO CARTAS VALUES (912,9,1,2);
INSERT INTO CARTAS VALUES (921,9,2,1);
INSERT INTO CARTAS VALUES (922,9,2,2);
INSERT INTO CARTAS VALUES (931,9,3,1);
INSERT INTO CARTAS VALUES (932,9,3,2);
INSERT INTO CARTAS VALUES (941,9,4,1);
INSERT INTO CARTAS VALUES (942,9,4,2);
INSERT INTO CARTAS VALUES (1011,10,1,1);
INSERT INTO CARTAS VALUES (1012,10,1,2);
INSERT INTO CARTAS VALUES (1021,10,2,1);
INSERT INTO CARTAS VALUES (1022,10,2,2);
INSERT INTO CARTAS VALUES (1031,10,3,1);
INSERT INTO CARTAS VALUES (1032,10,3,2);
INSERT INTO CARTAS VALUES (1041,10,4,1);
INSERT INTO CARTAS VALUES (1042,10,4,2);
INSERT INTO CARTAS VALUES (1111,11,1,1);
INSERT INTO CARTAS VALUES (1112,11,1,2);
INSERT INTO CARTAS VALUES (1121,11,2,1);
INSERT INTO CARTAS VALUES (1122,11,2,2);
INSERT INTO CARTAS VALUES (1131,11,3,1);
INSERT INTO CARTAS VALUES (1132,11,3,2);
INSERT INTO CARTAS VALUES (1141,11,4,1);
INSERT INTO CARTAS VALUES (1142,11,4,2);
INSERT INTO CARTAS VALUES (1211,12,1,1);
INSERT INTO CARTAS VALUES (1212,12,1,2);
INSERT INTO CARTAS VALUES (1221,12,2,1);
INSERT INTO CARTAS VALUES (1222,12,2,2);
INSERT INTO CARTAS VALUES (1231,12,3,1);
INSERT INTO CARTAS VALUES (1232,12,3,2);
INSERT INTO CARTAS VALUES (1241,12,4,1);
INSERT INTO CARTAS VALUES (1242,12,4,2);
INSERT INTO CARTAS VALUES (1311,13,1,1);
INSERT INTO CARTAS VALUES (1312,13,1,2);
INSERT INTO CARTAS VALUES (1321,13,2,1);
INSERT INTO CARTAS VALUES (1322,13,2,2);
INSERT INTO CARTAS VALUES (1331,13,3,1);
INSERT INTO CARTAS VALUES (1332,13,3,2);
INSERT INTO CARTAS VALUES (1341,13,4,1);
INSERT INTO CARTAS VALUES (1342,13,4,2);