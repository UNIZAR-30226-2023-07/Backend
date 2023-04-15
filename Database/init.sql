CREATE TABLE JUGADORES (
	nombre text NOT NULL,
    contra text NOT NULL,
	foto int NOT NULL,
    descrp text,
    pjugadas integer NOT NULL,
    pganadas integer NOT NULL,
	puntos integer NOT NULL,
    email text NOT NULL,
	codigo text,
	UNIQUE (email),
    PRIMARY KEY (codigo)
);

CREATE TABLE AMISTAD (
	id serial,
    estado text NOT NULL CHECK(estado = 'esp_confirmacion' OR estado = 'confirmada' OR estado = 'pendiente'),
    usr1 text REFERENCES JUGADORES (codigo),
    usr2 text REFERENCES JUGADORES (codigo),
	UNIQUE(usr1, usr2),
    PRIMARY KEY (id)
);

CREATE TABLE PARTIDAS (
	clave text,
	creador text REFERENCES JUGADORES (codigo),
	tipo text NOT NULL CHECK(tipo = 'amistosa' OR tipo = 'torneo'), 
--	fecha date NOT NULL,
	estado text NOT NULL CHECK(estado = 'terminada' OR estado = 'pausada' OR estado = 'iniciada' OR estado = 'creando'),
	PRIMARY KEY (clave)
);

CREATE TABLE PARTICIPAR (
	id serial,
	partida text REFERENCES PARTIDAS (clave),
	jugador text REFERENCES JUGADORES (codigo),
	puntos_resultado integer NOT NULL,
	UNIQUE (partida, jugador),
	PRIMARY KEY (id)
);

CREATE TABLE MENSAJES (
	id serial, 
	jug_emi text REFERENCES JUGADORES (codigo),
	jug_rcp text REFERENCES JUGADORES (codigo),
	timestamp timestamp NOT NULL,
	contenido text NOT NULL,
	leido integer NOT NULL,
	PRIMARY KEY(id)
);

CREATE TABLE CARTAS(
    id serial,
	numero integer NOT NULL,
    palo text NOT NULL,
	UNIQUE (numero, palo),
    PRIMARY KEY (id)
);	

insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, email, codigo) values ('Max', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 7, 'Aliquam sit amet diam in magna bibendum imperdiet. Nullam orci pede, venenatis non, sodales sed, tincidunt eu, felis. Fusce posuere felis sed lacus.', 9, 14, 1000, 'max@gmail.com', 1);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, email, codigo) values ('Eugenio', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 1, 'Proin at turpis a pede posuere nonummy. Integer non velit.', 12, 87, 900, 'eu@gmail.com', 2);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, email, codigo) values ('Lauren', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 2, 'Morbi non lectus.', 94, 88, 800, 'lclipson2@trellian.com', 3);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, email, codigo) values ('Shannon', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 9, 'Cras non velit nec nisi vulputate nonummy. Maecenas tincidunt lacus at velit.', 85, 32, 700, 'shan@gmail.com', 4);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, email, codigo) values ('Valentino', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 9, 'Donec vitae nisi. Nam ultrices, libero non mattis pulvinar, nulla,', 76, 88, 600, 'val@gmail.com', 5);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, email, codigo) values ('Ernest', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 2, 'Donec posuere metus vitae ipsum. Aliquam non mauris. Morbi non lectus.', 52, 90, 500, 'erni@gmail.com', 6);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, email, codigo) values ('Frederic', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 8, 'Fusce posuere felis sed lacus. Morbi sem mauris, laoreet ut, rhoncus aliquet, pulvinar sed, nisl.', 25, 31, 400, 'fredy@gmail.com', 7);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, email, codigo) values ('Edeline', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 10, 'In hac habitasse platea dictumst. Morbi vestibulum, velit id pretium iaculis, diam erat fermentum justo, nec condimentum neque sapien placerat ante.', 95, 75, 300, 'edy@gmail.com', 8);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, email, codigo) values ('Niki', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 3, 'Nullam varius.', 29, 23, 200, 'minaj@gmail.com', 9);
insert into JUGADORES (nombre, contra, foto, descrp, pjugadas, pganadas, email, codigo) values ('Luce', '03ac674216f3e15c761ee1a5e255f067953623c8b388b4459e13f978d7c846f4', 2, 'Fusce consequat. Nulla nisl. Nunc nisl.', 1, 99, 100, 'luce@gamil.com', 10);

insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 2);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 3);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 4);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 5);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 6);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 7);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 8);
insert into amistad (estado, usr1, usr2 ) values ('confirmada', 1, 9);
