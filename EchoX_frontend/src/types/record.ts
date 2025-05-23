export interface Tag {
  id: number;
  name: string;
  color: string;
  user_id: number;
  created_at: string;
  updated_at: string;
}

export interface Record {
  id: number;
  title: string;
  content: string;
  type: 'text' | 'image' | 'video';
  media_url?: string;
  user_id: number;
  tags: Tag[];
  created_at: string;
  updated_at: string;
}

export interface CreateRecordDto {
  title: string;
  content: string;
  type: 'text' | 'image' | 'video';
  media_url?: string;
  tag_ids?: number[];
}

export interface UpdateRecordDto {
  title?: string;
  content?: string;
  type?: 'text' | 'image' | 'video';
  media_url?: string;
  tag_ids?: number[];
} 