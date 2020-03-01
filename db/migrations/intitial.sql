CREATE TABLE ip_white_list (
   id serial PRIMARY KEY, 
   network CIDR NOT NULL
);

CREATE TABLE ip_black_list (
    id serial PRIMARY KEY, 
    network CIDR NOT NULL
);

