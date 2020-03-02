use std::time;

use actix::prelude::{Message, Recipient};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

use crate::auth::permissions::AuthLevel;

// Struct sent to users
#[derive(Message)]
#[rtype(result = "()")]
pub struct Broadcast(pub String);

#[derive(Debug, Message, Serialize, Deserialize)]
#[rtype(result = "()")]
pub struct ChatMessage {
    pub username: String,
    pub date: time::SystemTime,
    pub id: Uuid,
    pub text: String,
}

#[derive(Debug, Message)]
#[rtype(usize)]
pub struct Connect {
    pub addr: Recipient<Broadcast>,
}

#[derive(Debug, Message)]
#[rtype(result = "()")]
pub struct Disconnect {
    pub id: usize,
}

#[derive(Debug, Message)]
#[rtype(result = "()")]
pub struct Command {
    pub from: String,
    pub name: String,
    pub args: Vec<String>,
}

