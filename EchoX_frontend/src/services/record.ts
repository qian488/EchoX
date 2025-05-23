import axios from 'axios';
import type { Record, CreateRecordDto, UpdateRecordDto } from '../types/record';
import { useUserStore } from '../stores/user';

const API_URL = 'http://localhost:8080/api/v1';

export const recordService = {
  async getRecords(): Promise<Record[]> {
    const userStore = useUserStore();
    const response = await axios.get(`${API_URL}/records`, {
      headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    });
    return response.data;
  },

  async getRecord(id: number): Promise<Record> {
    const userStore = useUserStore();
    const response = await axios.get(`${API_URL}/records/${id}`, {
      headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    });
    return response.data;
  },

  async createRecord(record: CreateRecordDto): Promise<Record> {
    const userStore = useUserStore();
    const response = await axios.post(`${API_URL}/records`, record, {
      headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    });
    return response.data;
  },

  async updateRecord(id: number, record: UpdateRecordDto): Promise<Record> {
    const userStore = useUserStore();
    const response = await axios.put(`${API_URL}/records/${id}`, record, {
      headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    });
    return response.data;
  },

  async deleteRecord(id: number): Promise<void> {
    const userStore = useUserStore();
    await axios.delete(`${API_URL}/records/${id}`, {
      headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    });
  }
}; 