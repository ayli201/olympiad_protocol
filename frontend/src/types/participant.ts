import { type Task } from "./task";
import { type Timestamps } from "./timestamps";

export interface Participant extends Timestamps {
  id: number;
  fullName: string;
  cipher: number;
  organizationId: number;
  group: string;
  tasks: Task[];
  total: number;
  percent: number;
  rating: string;
  status: string;
}
