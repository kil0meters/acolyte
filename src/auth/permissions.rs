// only i8/i16/i32/i64 are supported by Insertable
pub type AuthLevel = i32;

pub const ADMIN: AuthLevel = 0;
pub const MODERATOR: AuthLevel = 1;
pub const STANDARD: AuthLevel = 2;
pub const LOGGED_OUT: AuthLevel = 3;
pub const BANNED: AuthLevel = 4;

/// Lower values give higher permissions
///
/// ```
/// const ADMIN      = 0;
/// const MODERATOR  = 1;
/// const STANDARD   = 2;
/// const LOGGED_OUT = 3;
/// const BANNED     = 4;
/// ```
///
/// We cannot use an enum because it doesn't integrate with Diesel well.
pub fn check_auth_level(test_value: AuthLevel, minimum_permission: AuthLevel) -> bool {
    test_value <= minimum_permission
}

pub trait Permission {
    fn at_least(&self, minimum_permission: AuthLevel) -> bool;
}

impl Permission for AuthLevel {
    fn at_least(&self, minimum_permission: AuthLevel) -> bool {
        self <= &minimum_permission
    }
}
