import { Show } from "@refinedev/antd";
import { useShow } from "@refinedev/core";
import { Typography, Tag, Space } from "antd";

const { Title, Text } = Typography;

export const ProjectShow = () => {
  const { queryResult } = useShow({
    resource: "projects",
  });

  const { data, isLoading } = queryResult;
  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Title level={5}>Title</Title>
      <Text>{record?.title}</Text>

      <Title level={5} style={{ marginTop: 16 }}>
        Slug
      </Title>
      <Text>{record?.slug}</Text>

      <Title level={5} style={{ marginTop: 16 }}>
        Category
      </Title>
      <Tag>{record?.category}</Tag>

      <Title level={5} style={{ marginTop: 16 }}>
        Status
      </Title>
      <Tag>{record?.status}</Tag>

      <Title level={5} style={{ marginTop: 16 }}>
        Featured
      </Title>
      <Text>{record?.featured ? "Yes" : "No"}</Text>

      {record?.short_description && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Short Description
          </Title>
          <Text>{record.short_description}</Text>
        </>
      )}

      {record?.full_description && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Full Description
          </Title>
          <Text>{record.full_description}</Text>
        </>
      )}

      {record?.cover_image_url && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Cover Image
          </Title>
          <img
            src={record.cover_image_url}
            alt={record.title}
            style={{ maxWidth: "100%", maxHeight: 400 }}
          />
        </>
      )}

      {record?.client_name && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Client
          </Title>
          <Text>{record.client_name}</Text>
        </>
      )}

      {record?.project_url && (
        <>
          <Title level={5} style={{ marginTop: 16 }}>
            Project URL
          </Title>
          <a href={record.project_url} target="_blank" rel="noopener noreferrer">
            {record.project_url}
          </a>
        </>
      )}
    </Show>
  );
};
