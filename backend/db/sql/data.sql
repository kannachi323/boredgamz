INSERT INTO users (name, email, password, is_admin)
VALUES(
    'Admin', 
    'admin@gmail.com', 
    crypt('poptropica911', gen_salt('bf')),
    true
);


