import { Subject } from "rxjs";

type AuthEvents = { type: "unauthorized" };

export const authSubject = new Subject<AuthEvents>();
