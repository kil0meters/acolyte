use super::message_types::ChatMessage;
use super::session::Client;
use crate::auth::permissions::{self, AuthLevel, Permission};

#[derive(Debug, Clone)]
pub struct Command {
    pub name: String,
    pub description: String,
    pub minimum_permissions: AuthLevel,
}

impl Command {
    pub fn to_string(&self) -> String {
        return format!("/{} {}", self.name, self.description);
    }
}

pub fn execute_command(user: Client, text: String) -> ChatMessage {
    let arguments: Vec<&str> = text.split(" ").collect();

    // this is probably a poor architecture for this, but it's usable for now
    if user.auth_level.at_least(permissions::MODERATOR) {
        match arguments[0] {
            "/addcommand" | "/ac" => {
                return ChatMessage::new("System", "Added command");
            }

            _ => {
                return ChatMessage::new("System", "ERROR: Unknown command");
            }
        }
    } else {
        return ChatMessage::new(
            "System",
            "ERROR: You don't have the required permissions to use that command",
        );
    }
}
