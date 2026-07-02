export class AppError extends Error {
  constructor(
    message: string,
    public userTitle: string,
    public userDescription: string,
    public actionLabel?: string,
  ) {
    super(message);
    this.name = "AppError";
  }
}

export class UserFacingError extends Error {
  constructor(
    public userTitle: string,
    public userDescription: string,
    public actionLabel?: string,
  ) {
    super(`${userTitle}: ${userDescription}`);
    this.name = "UserFacingError";
  }
}
