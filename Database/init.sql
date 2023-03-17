CREATE TABLE JUGADORES (
	nombre text NOT NULL,
    contra text NOT NULL,
	foto int NOT NULL,
    descrp text,
    pjugadas integer NOT NULL,
    pganadas integer NOT NULL,
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
	fecha date NOT NULL,
	estado text NOT NULL CHECK(estado = 'terminada' OR estado = 'pausada' OR estado = 'iniciada'),
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

CREATE TABLE CARTAS(
    id serial,
	numero integer NOT NULL,
    palo text NOT NULL,
	UNIQUE (numero, palo),
    PRIMARY KEY (id)
);	


