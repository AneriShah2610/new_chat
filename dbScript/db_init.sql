-- Db_Script


-- Drop table

-- DROP TABLE public.user_test
-- DROP TABLE public.chatroom_test
-- DROP TABLE public.members_test
-- DROP TABLE public.chatconversation

CREATE TABLE users
(
  id              INT       NOT NULL DEFAULT unique_rowid(),
  username        STRING    NOT NULL,
  first_name      STRING    NULL     DEFAULT NULL,
  last_name       STRING    NULL     DEFAULT NULL,
  email           STRING    NOT NULL,
  contact         STRING    NULL     DEFAULT NULL,
  bio             STRING    NULL     DEFAULT NULL,
  profile_picture STRING    NULL     DEFAULT NULL,
  created_at      TIMESTAMP NULL,
  updated_at      TIMESTAMP NULL     DEFAULT NULL,
  CONSTRAINT "primary" PRIMARY KEY (id ASC),
  UNIQUE INDEX chat_users_email_key (email ASC),
  FAMILY          "primary"(id, username, first_name, last_name, email, contact, bio, profile_picture, created_at, updated_at)
);

CREATE TABLE chatrooms
(
  id            INT       NOT NULL DEFAULT unique_rowid(),
  creator_id    INT       NOT NULL,
  chatroom_name STRING    NULL     DEFAULT NULL,
  chatroom_type STRING    NOT NULL,
  created_at    TIMESTAMP NOT NULL,
  updated_by    INT       NULL     DEFAULT NULL,
  updated_at    TIMESTAMP NULL     DEFAULT NULL,
  deleted_at    TIMESTAMP NULL     DEFAULT NULL,
  hashkey       STRING    NULL     DEFAULT NULL,
  CONSTRAINT "primary" PRIMARY KEY (id ASC),
  CONSTRAINT fk_creator_id_ref_users FOREIGN KEY (creator_id) REFERENCES users (id),
  INDEX         chatroom_auto_index_fk_creator_id_ref_users(creator_id ASC),
  CONSTRAINT fk_updated_by_ref_users FOREIGN KEY (updated_by) REFERENCES users (id),
  INDEX         chatroom_auto_index_fk_updated_by_ref_users(updated_by ASC),
  UNIQUE INDEX chatrooms_hashkey_key (hashkey ASC),
  UNIQUE INDEX creator_id_chatroom_name (creator_id ASC, chatroom_name ASC),
  FAMILY        "primary"(id, creator_id, chatroom_name, chatroom_type, created_at, updated_by, updated_at, deleted_at, hashkey)
);

CREATE TABLE members
(
  id          INT       NOT NULL DEFAULT unique_rowid(),
  chatroom_id INT       NOT NULL,
  member_id   INT       NOT NULL,
  joined_at   TIMESTAMP NOT NULL,
  deleted_at  TIMESTAMP NULL,
  delete_flag INT       NULL     DEFAULT 0:::INT,
  CONSTRAINT "primary" PRIMARY KEY (id ASC),
  CONSTRAINT fk_chatroom_id_ref_chatrooms FOREIGN KEY (chatroom_id) REFERENCES chatrooms (id),
  INDEX       members_auto_index_fk_chatroom_id_ref_chatrooms(chatroom_id ASC),
  CONSTRAINT fk_member_id_ref_users FOREIGN KEY (member_id) REFERENCES users (id) ON DELETE CASCADE,
  INDEX       members_auto_index_fk_member_id_ref_users(member_id ASC),
  UNIQUE INDEX chatroom_id_member_id_unique (chatroom_id ASC, member_id ASC),
  FAMILY      "primary"(id, chatroom_id, member_id, joined_at, deleted_at, delete_flag)
);


CREATE TABLE chatconversation
(
  id                INT       NOT NULL DEFAULT unique_rowid(),
  chatroom_id       INT       NOT NULL,
  sender_id         INT       NOT NULL,
  message           STRING    NOT NULL,
  message_type      STRING    NOT NULL,
  message_status    STRING    NOT NULL,
  message_parent_id INT       NULL,
  created_at        TIMESTAMP NOT NULL,
  updatedat         TIMESTAMP NULL     DEFAULT NULL,
  CONSTRAINT "primary" PRIMARY KEY (id ASC),
  CONSTRAINT fk_chatroom_id_ref_chatrooms FOREIGN KEY (chatroom_id) REFERENCES chatrooms (id) ON DELETE CASCADE,
  INDEX             chatconversation_auto_index_fk_chatroom_id_ref_chatrooms(chatroom_id ASC),
  CONSTRAINT fk_sender_id_ref_users FOREIGN KEY (sender_id) REFERENCES users (id),
  INDEX             chatconversation_auto_index_fk_sender_id_ref_users(sender_id ASC),
  CONSTRAINT fk_message_parent_id_ref_chatconversation FOREIGN KEY (message_parent_id) REFERENCES chatconversation (id),
  INDEX             chatconversation_auto_index_fk_message_parent_id_ref_chatconversation(message_parent_id ASC),
  FAMILY            "primary"(id, chatroom_id, sender_id, message, message_type, message_status, message_parent_id, created_at, updatedat)
);

