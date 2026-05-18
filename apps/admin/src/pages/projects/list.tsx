import { List, useTable, EditButton, ShowButton, DeleteButton } from "@refinedev/antd";
import { Table, Space, Tag } from "antd";

export const ProjectList = () => {
  const { tableProps } = useTable({
    resource: "projects",
    syncWithLocation: true,
  });

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="title" title="Title" />
        <Table.Column dataIndex="slug" title="Slug" />
        <Table.Column
          dataIndex="category"
          title="Category"
          render={(value) => {
            const colors: Record<string, string> = {
              web: "blue",
              logiciel: "green",
              ia: "purple",
              conseil: "orange",
            };
            return <Tag color={colors[value] || "default"}>{value}</Tag>;
          }}
        />
        <Table.Column
          dataIndex="status"
          title="Status"
          render={(value) => {
            const colors: Record<string, string> = {
              draft: "default",
              published: "success",
              archived: "error",
            };
            return <Tag color={colors[value] || "default"}>{value}</Tag>;
          }}
        />
        <Table.Column
          dataIndex="featured"
          title="Featured"
          render={(value) => (value ? "Yes" : "No")}
        />
        <Table.Column
          title="Actions"
          render={(_, record: any) => (
            <Space>
              <EditButton hideText size="small" recordItemId={record.id} />
              <ShowButton hideText size="small" recordItemId={record.id} />
              <DeleteButton hideText size="small" recordItemId={record.id} />
            </Space>
          )}
        />
      </Table>
    </List>
  );
};
