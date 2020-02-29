use std::time;

use actix::prelude::{Message, Recipient};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

pub enum AuthLevel {
    Standard,
    Moderator,
    Admin,
}

// Struct sent to users
#[derive(Message)]
#[rtype(result = "()")]
pub struct Broadcast(pub String);

#[derive(Message, Serialize, Deserialize)]
#[rtype(result = "()")]
pub struct ChatMessage {
    pub username: String,
    pub date: time::SystemTime,
    pub id: Uuid,
    pub text: String,
}

#[derive(Message)]
#[rtype(usize)]
pub struct Connect {
    pub addr: Recipient<Broadcast>,
}

#[derive(Message)]
#[rtype(result = "()")]
pub struct Disconnect {
    pub id: usize,
}

#[derive(Message)]
#[rtype(result = "()")]
pub struct Join {
    pub id: usize,
    pub username: String,
    pub auth_level: AuthLevel,
}

#[derive(Message)]
#[rtype(result = "()")]
pub struct Command {
    pub from: String,
    pub name: String,
    pub args: Vec<String>,
}

