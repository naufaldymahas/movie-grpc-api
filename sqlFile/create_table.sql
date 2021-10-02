create table movie_logs (
	id BIGINT not null primary key auto_increment,
	created_at datetime default (CURRENT_TIMESTAMP),
	event_type varchar(255),
	params text
);