UPDATE products p
SET quantity = quantity - sold_data.total_sold
FROM (
    SELECT product_id, SUM(quantity) as total_sold
    FROM order_items
    GROUP BY product_id
) AS sold_data
WHERE p.id = sold_data.product_id;

######################################################################################################
# Correcting the stock quantities for products 
# use the command sudo docker exec -i ecom-postgres psql -U postgres -d ecom < quantity_correction.sql
######################################################################################################