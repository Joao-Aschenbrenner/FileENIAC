export class ApiError extends Error {
  constructor(public status: number, public url: string, message: string) {
    super(message);
    this.name = "ApiError";
  }
  isUnauthorized(): boolean {
    return this.status === 401;
  }
  isForbidden(): boolean {
    return this.status === 403;
  }
  isTimeout(): boolean {
    return false;
  }
}
