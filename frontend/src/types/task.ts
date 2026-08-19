import { type Timestamps } from "./timestamps";

export interface Task extends Timestamps {
  id: number;
  value: number;
  number: number;
}
