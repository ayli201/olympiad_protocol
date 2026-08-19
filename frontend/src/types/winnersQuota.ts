import { type Timestamps } from "./timestamps";

export interface WinnersQuota extends Timestamps {
  id: number;
  minValue: number;
  maxValue: number;
  quota: number;
}
