import axios from 'axios';
import type { Tag } from '../types/record';
import { useUserStore } from '../stores/user';

const API_URL = 'http://localhost:8080/api/v1';

export const tagService = {
  async getTags(): Promise<Tag[]> {
    const userStore = useUserStore();
    const response = await axios.get(`${API_URL}/tags`, {
      headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    });
    return response.data;
  },

  async createTag(name: string, color: string): Promise<Tag> {
    const userStore = useUserStore();
    const response = await axios.post(`${API_URL}/tags`, { name, color }, {
       headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    });
    return response.data;
  },

  async updateTag(id: number, name: string, color: string): Promise<Tag> {
    const userStore = useUserStore();
    const response = await axios.put(`${API_URL}/tags/${id}`, { name, color }, {
       headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    });
    return response.data;
  },

  async deleteTag(id: number): Promise<void> {
    const userStore = useUserStore();
    await axios.delete(`${API_URL}/tags/${id}`, {
       headers: {
        Authorization: `Bearer ${userStore.token}`
      }
    });
  }
}; 