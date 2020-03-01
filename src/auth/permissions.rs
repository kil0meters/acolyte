/// AuthLevel
/// Lower values give higher permissions
#[derive(Debug)]
pub enum AuthLevel {
    Admin = 0,
    Moderator = 1,
    Standard = 2,
    LoggedOut = 3,
    Banned = 4,
}

impl AuthLevel {
    pub fn is_at_least(self, minimum: AuthLevel) -> bool {
        self as u8 <= minimum as u8
    }
}
