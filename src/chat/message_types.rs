use std::time;

use actix::prelude::{Message, Recipient};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

use super::session::Client;

// Struct sent to users
#[derive(Message)]
#[rtype(result = "()")]
pub struct Broadcast(pub String);

/// Struct containing data for a command sent through a live chat
///
/// ```rust
/// pub client: Client,
/// pub date: time::SystemTime,
/// pub id: Uuid,
/// pub text: String,
/// ```
///
#[derive(Debug, Message)]
#[rtype(result = "()")]
pub struct ChatCommand {
    pub client: Client,
    pub date: time::SystemTime,
    pub id: Uuid,
    pub text: String,
}

impl ChatCommand {
    pub fn new(client: Client, text: &str) -> ChatCommand {
        return ChatCommand {
            client,
            id: Uuid::new_v4(),
            date: time::SystemTime::now(),
            text: text.to_owned(),
        };
    }
}

/// Struct containing data for a message sent through the live chat.
///
/// ```rust
/// pub username: String,
/// pub date: time::SystemTime,
/// pub id: Uuid,
/// pub text: String,
/// ```
#[derive(Debug, Message, Serialize, Deserialize)]
#[rtype(result = "()")]
pub struct ChatMessage {
    pub username: String,
    pub date: time::SystemTime,
    pub id: Uuid,
    pub text: String,
}

impl ChatMessage {
    pub fn new(username: &str, text: &str) -> ChatMessage {
        return ChatMessage {
            username: username.to_owned(),
            id: Uuid::new_v4(),
            date: time::SystemTime::now(),
            text: text.to_owned(),
        };
    }
}

#[derive(Debug, Message)]
#[rtype(usize)]
pub struct Connect {
    pub client: Client,
    pub addr: Recipient<Broadcast>,
}

#[derive(Debug, Message)]
#[rtype(result = "()")]
pub struct Disconnect {
    pub id: usize,
}

/// A struct representing a parsed command structure
/// (any chat message starting with "/")
///
/// ```rust
/// pub from: String,
/// pub name: String,
/// pub args: Vec<String>
/// ```
#[derive(Debug, Message)]
#[rtype(result = "()")]
pub struct Command {
    pub from: String,
    pub name: String,
    pub args: Vec<String>,
}
