INSERT INTO products (name, price_in_cents, quantity) VALUES
('Mechanical Keyboard', 9900, 10),
('Gaming Mouse', 4900, 25),
('UltraWide Monitor', 45000, 5),
('USB-C Hub', 2500, 50)
ON CONFLICT (id) DO NOTHING;

#######################################################################################
# ADDING some test data 
# use the command sudo docker exec -i ecom-postgres psql -U postgres -d ecom < seed.sql
#######################################################################################