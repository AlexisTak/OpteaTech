import { Create, useForm } from "@refinedev/antd";
import { Form, Input, Select, Switch } from "antd";

const { TextArea } = Input;

export const ProjectCreate = () => {
  const { formProps, saveButtonProps } = useForm({
    resource: "projects",
  });

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Title"
          name="title"
          rules={[{ required: true, message: "Please enter title" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="Slug"
          name="slug"
          rules={[{ required: true, message: "Please enter slug" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item
          label="Category"
          name="category"
          rules={[{ required: true, message: "Please select category" }]}
        >
          <Select>
            <Select.Option value="web">Web</Select.Option>
            <Select.Option value="logiciel">Logiciel</Select.Option>
            <Select.Option value="ia">IA</Select.Option>
            <Select.Option value="conseil">Conseil</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item
          label="Status"
          name="status"
          initialValue="draft"
        >
          <Select>
            <Select.Option value="draft">Draft</Select.Option>
            <Select.Option value="published">Published</Select.Option>
            <Select.Option value="archived">Archived</Select.Option>
          </Select>
        </Form.Item>

        <Form.Item label="Short Description" name="short_description">
          <TextArea rows={3} />
        </Form.Item>

        <Form.Item label="Full Description" name="full_description">
          <TextArea rows={6} />
        </Form.Item>

        <Form.Item label="Cover Image URL" name="cover_image_url">
          <Input />
        </Form.Item>

        <Form.Item label="Project URL" name="project_url">
          <Input />
        </Form.Item>

        <Form.Item label="GitHub URL" name="github_url">
          <Input />
        </Form.Item>

        <Form.Item label="Client Name" name="client_name">
          <Input />
        </Form.Item>

        <Form.Item label="Featured" name="featured" valuePropName="checked">
          <Switch />
        </Form.Item>
      </Form>
    </Create>
  );
};
