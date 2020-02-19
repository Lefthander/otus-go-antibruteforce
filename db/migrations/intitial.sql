CREATE TABLE ip_white_list (
   id serial PRIMARY KEY, 
   network CIDR NOT NULL
);

CREATE TABLE ip_black_list (
    id serial PRIMARY KEY, 
    network CIDR NOT NULL
);
CREATE INDEX indx_white ON ip_white_list USING GIST(network);
CREATE INDEX indx_black ON ip_black_list USING GIST(network);  

