import axios from 'axios';

const API_URL = 'http://localhost:8080/api/v1';

/** 上传 / 查询媒体文件时的响应结构（url 为相对路径，需拼接后端地址展示） */
export interface MediaResponse {
  id: number;
  url: string;
  filename: string;
  mime_type: string;
  size: number;
  type: string;
}

const getAuthHeaders = () => ({
  Authorization: `Bearer ${localStorage.getItem('token') || ''}`
});

export const mediaService = {
  /**
   * 上传文件。FormData 的 Content-Type 与 boundary 由浏览器自动设置，不要手动指定。
   * onProgress 为可选的上传进度回调（0-100）。
   */
  async upload(
    file: File,
    type: 'text' | 'image' | 'video',
    onProgress?: (percent: number) => void
  ): Promise<MediaResponse> {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('type', type);

    const response = await axios.post<MediaResponse>(`${API_URL}/upload`, formData, {
      headers: getAuthHeaders(),
      onUploadProgress: onProgress
        ? (event) => {
            if (event.total && event.total > 0) {
              onProgress(Math.round((event.loaded * 100) / event.total));
            }
          }
        : undefined
    });
    return response.data;
  },

  /** 查询单个媒体文件信息 */
  async getMedia(id: number): Promise<MediaResponse> {
    const response = await axios.get<MediaResponse>(`${API_URL}/media/${id}`, {
      headers: getAuthHeaders()
    });
    return response.data;
  },

  /** 删除媒体文件（同时删除数据库记录与物理文件） */
  async deleteMedia(id: number): Promise<void> {
    await axios.delete(`${API_URL}/media/${id}`, {
      headers: getAuthHeaders()
    });
  }
};
