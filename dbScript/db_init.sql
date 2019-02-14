-- Db_Script


-- Drop table

-- DROP TABLE public.user_test
-- DROP TABLE public.chatroom_test
-- DROP TABLE public.members_test
-- DROP TABLE public.chatconversation


CREATE TABLE user_test (
	id INT NOT NULL DEFAULT unique_rowid(),
	"name" STRING NOT NULL,
	email STRING NULL DEFAULT NULL,
	contact STRING NOT NULL,
	profile_picture STRING NULL DEFAULT NULL,
	bio STRING NULL DEFAULT NULL,
	createdat TIMESTAMP WITH TIME ZONE NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	UNIQUE INDEX user_test_name_key ("name" ASC),
	UNIQUE INDEX user_test_contact_key (contact ASC),
	FAMILY "primary" (id, "name", email, contact, profile_picture, bio, createdat)
);


CREATE TABLE chatroom_test (
	id INT NOT NULL DEFAULT unique_rowid(),
	creator_id INT NOT NULL,
	chatroom_name STRING NULL DEFAULT NULL,
	chatroom_type STRING NOT NULL,
	createat TIMESTAMP WITH TIME ZONE NOT NULL,
	updateby INT NULL,
	updateat TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
	deleteat TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
	hashkey STRING NULL DEFAULT NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	CONSTRAINT fk_creator_id_ref_user_test FOREIGN KEY (creator_id) REFERENCES user_test (id) ON DELETE CASCADE,
	INDEX chatroom_test_auto_index_fk_creator_id_ref_user_test (creator_id ASC),
	CONSTRAINT fk_updateby_ref_user_test FOREIGN KEY (updateby) REFERENCES user_test (id),
	INDEX chatroom_test_auto_index_fk_updateby_ref_user_test (updateby ASC),
	UNIQUE INDEX chatroom_test_hashkey_key (hashkey ASC),
	FAMILY "primary" (id, creator_id, chatroom_name, chatroom_type, createat, updateby, updateat, deleteat, hashkey)
);

CREATE TABLE members_test (
	id INT NOT NULL DEFAULT unique_rowid(),
	chatroom_id INT NOT NULL,
	member_id INT NOT NULL,
	joinat TIMESTAMP WITH TIME ZONE NOT NULL,
	deleteat TIMESTAMP WITH TIME ZONE NULL DEFAULT NULL,
	flag BOOL NULL DEFAULT false,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	CONSTRAINT fk_chatroom_id_ref_chatroom_test FOREIGN KEY (chatroom_id) REFERENCES chatroom_test (id),
	INDEX members_test_auto_index_fk_chatroom_id_ref_chatroom_test (chatroom_id ASC),
	CONSTRAINT fk_member_id_ref_user_test FOREIGN KEY (member_id) REFERENCES user_test (id),
	INDEX members_test_auto_index_fk_member_id_ref_user_test (member_id ASC),
	FAMILY "primary" (id, chatroom_id, member_id, joinat, deleteat, flag)
);

CREATE TABLE chatconversation (
	id INT NOT NULL DEFAULT unique_rowid(),
	chatroom_id INT NOT NULL,
	sender_id INT NOT NULL,
	message STRING NOT NULL,
	message_type STRING NOT NULL,
	message_parent_id INT NULL DEFAULT NULL,
	message_status STRING NULL DEFAULT 'SEND':::STRING,
	createat TIMESTAMP NOT NULL,
	updateat TIMESTAMP NULL DEFAULT NULL,
	deleteat TIMESTAMP NULL DEFAULT NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	CONSTRAINT fk_chatroom_id_ref_chatroom_test FOREIGN KEY (chatroom_id) REFERENCES chatroom_test (id),
	INDEX chatconversation_auto_index_fk_chatroom_id_ref_chatroom_test (chatroom_id ASC),
	CONSTRAINT fk_sender_id_ref_user_test FOREIGN KEY (sender_id) REFERENCES user_test (id),
	INDEX chatconversation_auto_index_fk_sender_id_ref_user_test (sender_id ASC),
	CONSTRAINT fk_message_parent_id_ref_chatconversation FOREIGN KEY (message_parent_id) REFERENCES chatconversation (id),
	INDEX chatconversation_auto_index_fk_message_parent_id_ref_chatconversation (message_parent_id ASC),
	FAMILY "primary" (id, chatroom_id, sender_id, message, message_type, message_parent_id, message_status, createat, updateat, deleteat)
);

