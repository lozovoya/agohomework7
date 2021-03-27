CREATE TABLE flights (
                         id BIGSERIAL PRIMARY KEY,
                         iata_from TEXT NOT NULL,
                         iata_to TEXT NOT NULL,
                         departure_date DATE NOT NULL,
                         departure_time TIME WITH TIME ZONE NOT NULL,
                         arrival_date DATE NOT NULL,
                         arrival_time TIME WITH TIME ZONE NOT NULL,
                         duration INT NOT NULL,
                         price INT NOT NULL,
                         departure TIMESTAMP WITH TIME ZONE
);

