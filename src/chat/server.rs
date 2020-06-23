use std::collections::HashMap;
use std::sync::RwLock;

use actix::prelude::{Context, Handler, Recipient};
use actix::Actor;
use rand::{self, rngs::ThreadRng, Rng};
use serde_json;

use super::commands::{self, Command};
use super::message_types::*;
use super::session::Client;
use crate::auth::permissions::{self, AuthLevel, Permission};

#[derive(Debug)]
pub struct Server {
    // The key is the username of the session, while the
    sessions: HashMap<usize, (Client, Recipient<Broadcast>)>,
    commands: RwLock<Vec<Command>>,
    rng: ThreadRng,
}

impl Default for Server {
    fn default() -> Server {
        Server {
            sessions: HashMap::new(),

            // commands used to generate auto completion
            commands: RwLock::new(vec![Command {
                name: "addcommand".to_string(),
                description: "[name] [output...]".to_string(),
                minimum_permissions: permissions::MODERATOR,
            }]),
            rng: rand::thread_rng(),
        }
    }
}

impl Server {
    // fn find_session_from_username(&self, username: String) -> Option<Recipient<Broadcast>> {
    //     return self.sessions.get(username);
    // }

    fn broadcast_message(&self, message: String) {
        for (_username, addr) in self.sessions.values() {
            let _ = addr.do_send(Broadcast(message.to_owned()));
        }
    }

    fn send_to_user(&self, username: String, msg: ChatMessage) {
        let msg_str = serde_json::to_string(&msg).unwrap();

        for (client, addr) in self.sessions.values() {
            if client.username == username {
                let _ = addr.do_send(Broadcast(msg_str.clone()));
            }
        }
    }
}

impl Actor for Server {
    type Context = Context<Self>;
}

impl Handler<ChatCommand> for Server {
    type Result = ();

    fn handle(&mut self, cmd: ChatCommand, _: &mut Context<Self>) {
        self.send_to_user(
            cmd.client.username.clone(),
            commands::execute_command(cmd.client, cmd.text),
        );
    }
}

impl Handler<ChatMessage> for Server {
    type Result = ();

    fn handle(&mut self, msg: ChatMessage, _: &mut Context<Self>) {
        let m = serde_json::to_string(&msg).unwrap();
        self.broadcast_message(m);
    }
}

impl Handler<Connect> for Server {
    type Result = usize;

    fn handle(&mut self, msg: Connect, _: &mut Context<Self>) -> Self::Result {
        let id = self.rng.gen::<usize>();
        self.sessions.insert(id, (msg.client, msg.addr));

        id
    }
}

impl Handler<Disconnect> for Server {
    type Result = ();

    fn handle(&mut self, msg: Disconnect, _: &mut Context<Self>) {
        self.sessions.remove(&msg.id);
    }
}
