export interface User {
  id: string;
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface CreateUserPayload {
  email: string;
  password: string;
  firstName?: string;
  lastName?: string;
}
