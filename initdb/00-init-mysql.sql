-- create table product
create table if not exists product (
    id serial primary key,
    name varchar(100),
    price float,
    description varchar(250),
    image varchar(250)
);

-- create table customer
create table if not exists customer (
    id varchar(150) primary key,
    name varchar(100),
    email varchar(250),
    password varchar(250)
);

-- create table order
create table if not exists orders (
    id varchar(250) primary key,
    id_customer varchar(250),
    date timestamp,
    status varchar(250),
    FOREIGN KEY (id_customer)references customer(id)
);

-- create table admin
create table if not exists admin (
    id varchar(150) primary key,
    name varchar(100),
    email varchar(250),
    password varchar(250)
);

-- create table product_order
create table if not exists product_order (
    id varchar(250) primary key,
    id_order varchar(250),
    id_product bigint UNSIGNED,
    quantity integer,
    FOREIGN KEY (id_order) references orders(id),
    FOREIGN KEY (id_product) references product(id)
);

-- query insert admin
-- password admin
INSERT INTO public.admin (id, "name", email, "password")
VALUES('3fd113c4-b522-4eec-acb0-d24785bc43c7', 'admin', 'admin@admin.admin', '$2a$14$LZ/K.8FPx7n2RygyDFmqAO.f3YTu7qF32RPlppPehmRuZeeYr/tFa');


-- query insert dummy data
INSERT INTO product (name, price, description, image) VALUES
('Wireless Bluetooth Earbuds', 49.99, 'Experience high-quality sound and wireless freedom with these sleek earbuds. Perfect for working out or commuting.', 'https://dummyurl.com/wireless-earbuds.jpg'),
('Smartphone Tripod', 24.99, 'Capture stunning photos and videos with ease using this versatile smartphone tripod. Includes adjustable legs and a remote control.', 'https://dummyurl.com/smartphone-tripod.jpg'),
('Fitness Tracker Watch', 89.99, 'Stay on top of your health and fitness goals with this advanced tracker watch. Includes heart rate monitor, GPS tracking, and more.', 'https://dummyurl.com/fitness-tracker.jpg'),
('Wireless Charger', 29.99, 'Charge your devices quickly and conveniently with this sleek wireless charger. Includes LED indicator lights and automatic shutoff.', 'https://dummyurl.com/wireless-charger.jpg'),
('Portable Bluetooth Speaker', 69.99, 'Take your music on the go with this powerful portable speaker. Includes built-in microphone for hands-free calling.', 'https://dummyurl.com/portable-speaker.jpg'),
('Wireless Gaming Mouse', 39.99, 'Dominate the competition with this high-performance gaming mouse. Features customizable RGB lighting and six programmable buttons.', 'https://dummyurl.com/gaming-mouse.jpg'),
('USB-C Hub', 49.99, 'Expand your device''s capabilities with this versatile USB-C hub. Includes multiple USB ports, HDMI output, and SD card reader.', 'https://dummyurl.com/usb-c-hub.jpg'),
('Smart Home Security Camera', 129.99, 'Keep an eye on your home from anywhere with this advanced security camera. Features motion detection, night vision, and two-way audio.', 'https://dummyurl.com/security-camera.jpg'),
('Bluetooth Noise-Canceling Headphones', 149.99, 'Experience immersive audio and total comfort with these noise-canceling headphones. Includes up to 30 hours of battery life.', 'https://dummyurl.com/noise-canceling-headphones.jpg'),
('Portable Power Bank', 34.99, 'Never run out of battery power again with this handy portable power bank. Includes multiple charging ports and LED indicator lights.', 'https://dummyurl.com/power-bank.jpg'),
('Smartwatch', 199.99, 'Stay connected and organized with this advanced smartwatch. Includes GPS tracking, heart rate monitor, and mobile notifications.', 'https://dummyurl.com/smartwatch.jpg'),
('Bluetooth Keyboard', 79.99, 'Increase productivity and comfort with this wireless keyboard. Features ergonomic design and long battery life.', 'https://dummyurl.com/bluetooth-keyboard.jpg'),
('Wireless Gaming Headset', 89.99, 'Immerse yourself in the game with this powerful wireless headset. Includes noise-cancellation technology and up to 15 hours of battery life.', 'https://dummyurl.com/gaming-headset.jpg'),
('Smart Wi-Fi Plug', 29.99, 'Control your devices remotely with this handy smart plug. Compatible with Alexa and Google Assistant for voice control.', 'https://dummyurl.com/wifi-plug.jpg');
