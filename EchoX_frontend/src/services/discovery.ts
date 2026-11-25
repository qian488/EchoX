import axios from 'axios';
import type { Record } from '../types/record';

const API_URL = 'http://localhost:8080/api/v1';

/** 时间轴分组：按天（date 为 YYYY-MM-DD）或按月（date 为 YYYY-MM） */
export interface TimelineGroup {
  date: string;
  count: number;
  records: Record[];
}

/** 搜索参数，全部可选 */
export interface SearchParams {
  q?: string;
  type?: string;
  tag?: string;
  start?: string;
  end?: string;
}

const getAuthHeaders = () => ({
  Authorization: `Bearer ${localStorage.getItem('token') || ''}`
});

export const discoveryService = {
  /** 获取当前用户的时间轴，按日期倒序 */
  async getTimeline(group: 'day' | 'month'): Promise<TimelineGroup[]> {
    const response = await axios.get<TimelineGroup[]>(`${API_URL}/timeline`, {
      params: { group },
      headers: getAuthHeaders()
    });
    return response.data;
  },

  /** 按关键词、类型、标签名、日期范围搜索当前用户的记录 */
  async search(params: SearchParams): Promise<Record[]> {
    const response = await axios.get<Record[]>(`${API_URL}/search`, {
      params,
      headers: getAuthHeaders()
    });
    return response.data;
  }
};
